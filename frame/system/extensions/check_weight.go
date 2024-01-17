package system

import (
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/frame/system"
	"github.com/LimeChain/gosemble/primitives/log"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

type CheckWeight primitives.Weight

func (_ CheckWeight) AdditionalSigned() (ok sc.Empty, err primitives.TransactionValidityError) {
	ok = sc.Empty{}
	return ok, err
}

func (_ CheckWeight) Validate(_who *primitives.Address32, _call *primitives.Call, info *primitives.DispatchInfo, length sc.Compact) (ok primitives.ValidTransaction, err primitives.TransactionValidityError) {
	return DoValidate(info, length)
}

func (_ CheckWeight) ValidateUnsigned(_call *primitives.Call, info *primitives.DispatchInfo, length sc.Compact) (ok primitives.ValidTransaction, err primitives.TransactionValidityError) {
	return DoValidate(info, length)
}

func (_ CheckWeight) PreDispatch(_who *primitives.Address32, _call *primitives.Call, info *primitives.DispatchInfo, length sc.Compact) (ok primitives.Pre, err primitives.TransactionValidityError) {
	_, err = DoPreDispatch(info, length)
	return ok, err
}

func (_ CheckWeight) PreDispatchUnsigned(_call *primitives.Call, info *primitives.DispatchInfo, length sc.Compact) (ok primitives.Pre, err primitives.TransactionValidityError) {
	_, err = DoPreDispatch(info, length)
	return ok, err
}

func (_ CheckWeight) PostDispatch(_pre sc.Option[primitives.Pre], info *primitives.DispatchInfo, postInfo *primitives.PostDispatchInfo, _length sc.Compact, _result *primitives.DispatchResult) (primitives.Pre, primitives.TransactionValidityError) {
	unspent := postInfo.CalcUnspent(info)
	if unspent.AnyGt(primitives.WeightZero()) {
		currentWeight := system.StorageGetBlockWeight()
		currentWeight.Reduce(unspent, info.Class)
		system.StorageSetBlockWeight(currentWeight)
	}
	return primitives.Pre{}, nil
}

// Do the validate checks. This can be applied to both signed and unsigned.
//
// It only checks that the block weight and length limit will not exceed.
func DoValidate(info *primitives.DispatchInfo, length sc.Compact) (ok primitives.ValidTransaction, err primitives.TransactionValidityError) {
	ok = primitives.DefaultValidTransaction()

	// ignore the next length. If they return `Ok`, then it is below the limit.
	_, err = checkBlockLength(info, length)
	if err != nil {
		return ok, err
	}

	// during validation we skip block limit check. Since the `validate_transaction`
	// call runs on an empty block anyway, by this we prevent `on_initialize` weight
	// consumption from causing false negatives.
	_, err = checkExtrinsicWeight(info)
	if err != nil {
		return ok, err
	}

	return ok, err
}

func DoPreDispatch(info *primitives.DispatchInfo, length sc.Compact) (ok primitives.ValidTransaction, err primitives.TransactionValidityError) {
	nextLength, err := checkBlockLength(info, length)
	if err != nil {
		return ok, err
	}

	nextWeight, err := checkBlockWeight(info)
	if err != nil {
		return ok, err
	}

	_, err = checkExtrinsicWeight(info)
	if err != nil {
		return ok, err
	}

	system.StorageSetAllExtrinsicsLen(nextLength)
	system.StorageSetBlockWeight(nextWeight)

	return ok, err
}

// Checks if the current extrinsic can fit into the block with respect to block length limits.
//
// Upon successes, it returns the new block length as a `Result`.
func checkBlockLength(info *primitives.DispatchInfo, length sc.Compact) (ok sc.U32, err primitives.TransactionValidityError) {
	lengthLimit := system.DefaultBlockLength()
	currentLen := system.StorageGetAllExtrinsicsLen()
	addedLen := sc.U32(sc.U128(length).ToBigInt().Uint64())

	nextLen := currentLen.SaturatingAdd(addedLen)

	var maxLimit sc.U32
	if info.Class.Is(primitives.DispatchClassNormal) {
		maxLimit = lengthLimit.Max.Normal
	} else if info.Class.Is(primitives.DispatchClassOperational) {
		maxLimit = lengthLimit.Max.Operational
	} else if info.Class.Is(primitives.DispatchClassMandatory) {
		maxLimit = lengthLimit.Max.Mandatory
	} else {
		log.Critical("invalid DispatchClass type in CheckBlockLength()")
	}

	if nextLen > maxLimit {
		err = primitives.NewTransactionValidityError(primitives.NewInvalidTransactionExhaustsResources())
	} else {
		ok = sc.U32(sc.ToCompact(nextLen).ToBigInt().Uint64())
	}

	return ok, err
}

// Checks if the current extrinsic can fit into the block with respect to block weight limits.
//
// Upon successes, it returns the new block weight as a `Result`.
func checkBlockWeight(info *primitives.DispatchInfo) (ok primitives.ConsumedWeight, err primitives.TransactionValidityError) {
	maximumWeight := system.DefaultBlockWeights()
	allWeight := system.StorageGetBlockWeight()
	return CalculateConsumedWeight(maximumWeight, allWeight, info)
}

// Checks if the current extrinsic does not exceed the maximum weight a single extrinsic
// with given `DispatchClass` can have.
func checkExtrinsicWeight(info *primitives.DispatchInfo) (ok sc.Empty, err primitives.TransactionValidityError) {
	max := system.DefaultBlockWeights().Get(info.Class).MaxExtrinsic

	if max.HasValue {
		if info.Weight.AnyGt(max.Value) {
			err = primitives.NewTransactionValidityError(primitives.NewInvalidTransactionExhaustsResources())
		} else {
			ok = sc.Empty{}
		}
	}

	return ok, err
}

func CalculateConsumedWeight(maximumWeight system.BlockWeights, allConsumedWeight primitives.ConsumedWeight, info *primitives.DispatchInfo) (ok primitives.ConsumedWeight, err primitives.TransactionValidityError) {
	limitPerClass := maximumWeight.Get(info.Class)
	extrinsicWeight := info.Weight.SaturatingAdd(limitPerClass.BaseExtrinsic)

	// add the weight. If class is unlimited, use saturating add instead of checked one.
	if !limitPerClass.MaxTotal.HasValue && !limitPerClass.Reserved.HasValue {
		allConsumedWeight.Accrue(extrinsicWeight, info.Class)
	} else {
		_, e := allConsumedWeight.CheckedAccrue(extrinsicWeight, info.Class)
		if e != nil {
			err = primitives.NewTransactionValidityError(primitives.NewInvalidTransactionExhaustsResources())
			return ok, err
		}
	}

	consumedPerClass := allConsumedWeight.Get(info.Class)

	// Check if we don't exceed per-class allowance
	if limitPerClass.MaxTotal.HasValue {
		max := limitPerClass.MaxTotal.Value
		if consumedPerClass.AnyGt(max) {
			err = primitives.NewTransactionValidityError(primitives.NewInvalidTransactionExhaustsResources())
			return ok, err
		}
	} else {
		// There is no `max_total` limit (`None`),
		// or we are below the limit.
	}

	// In cases total block weight is exceeded, we need to fall back
	// to `reserved` pool if there is any.
	if allConsumedWeight.Total().AnyGt(maximumWeight.MaxBlock) {
		if limitPerClass.Reserved.HasValue {
			reserved := limitPerClass.Reserved.Value
			if consumedPerClass.AnyGt(reserved) {
				err = primitives.NewTransactionValidityError(primitives.NewInvalidTransactionExhaustsResources())
				return ok, err
			}
		} else {
			// There is either no limit in reserved pool (`None`),
			// or we are below the limit.
		}
	}

	return allConsumedWeight, err
}
