package transaction_payment

import (
	"bytes"
	"math/big"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants"
	"github.com/LimeChain/gosemble/execution/types"
	"github.com/LimeChain/gosemble/frame/system"
	"github.com/LimeChain/gosemble/primitives/hashing"
	"github.com/LimeChain/gosemble/primitives/storage"
	primitives "github.com/LimeChain/gosemble/primitives/types"
	"github.com/LimeChain/gosemble/utils"
)

var DefaultMultiplierValue = sc.NewU128FromUint64(1)
var DefaultTip = sc.NewU128FromUint64(0)

// QueryInfo queries the data of an extrinsic.
// It takes two arguments:
// - dataPtr: Pointer to the data in the Wasm memory.
// - dataLen: Length of the data.
// which represent the SCALE-encoded extrinsic and its length.
// Returns a pointer-size of the SCALE-encoded weight, dispatch class and partial fee.
// [Specification](https://spec.polkadot.network/chap-runtime-api#sect-rte-transactionpaymentapi-query-info)
func QueryInfo(dataPtr int32, dataLen int32) int64 {
	b := utils.ToWasmMemorySlice(dataPtr, dataLen)
	buffer := bytes.NewBuffer(b)

	ext := types.DecodeUncheckedExtrinsic(buffer)
	length := sc.DecodeU32(buffer)

	dispatchInfo := primitives.GetDispatchInfo(ext.Function)

	partialFee := sc.NewU128FromUint64(0)
	if ext.IsSigned() {
		partialFee = computeFee(length, dispatchInfo, DefaultTip)
	}

	runtimeDispatchInfo := primitives.RuntimeDispatchInfo{
		Weight:     dispatchInfo.Weight,
		Class:      dispatchInfo.Class,
		PartialFee: partialFee,
	}

	return utils.BytesToOffsetAndSize(runtimeDispatchInfo.Bytes())
}

// QueryFeeDetails queries the detailed fee of an extrinsic.
// It takes two arguments:
// - dataPtr: Pointer to the data in the Wasm memory.
// - dataLen: Length of the data.
// which represent the SCALE-encoded extrinsic and its length.
// Returns a pointer-size of the SCALE-encoded detailed fee.
// [Specification](https://spec.polkadot.network/chap-runtime-api#sect-rte-transactionpaymentapi-query-fee-details)
func QueryFeeDetails(dataPtr int32, dataLen int32) int64 {
	b := utils.ToWasmMemorySlice(dataPtr, dataLen)
	buffer := bytes.NewBuffer(b)

	ext := types.DecodeUncheckedExtrinsic(buffer)
	length := sc.DecodeU32(buffer)

	dispatchInfo := primitives.GetDispatchInfo(ext.Function)

	var feeDetails primitives.FeeDetails
	if ext.IsSigned() {
		feeDetails = computeFeeDetails(length, dispatchInfo, DefaultTip)
	} else {
		feeDetails = primitives.FeeDetails{
			InclusionFee: sc.NewOption[primitives.InclusionFee](nil),
		}
	}

	return utils.BytesToOffsetAndSize(feeDetails.Bytes())
}

// QueryCallInfo queries the data of a dispatch call.
// It takes two arguments:
// - dataPtr: Pointer to the data in the Wasm memory.
// - dataLen: Length of the data.
// which represent the SCALE-encoded dispatch call and its length.
// Returns a pointer-size of the SCALE-encoded weight, dispatch class and partial fee.
// [Specification](https://spec.polkadot.network/chap-runtime-api#sect-rte-transactionpaymentcallapi-query-call-info)
func QueryCallInfo(dataPtr int32, dataLen int32) int64 {
	b := utils.ToWasmMemorySlice(dataPtr, dataLen)
	buffer := bytes.NewBuffer(b)

	call := types.DecodeCall(buffer)
	length := sc.DecodeU32(buffer)

	dispatchInfo := primitives.GetDispatchInfo(call)
	partialFee := computeFee(length, dispatchInfo, DefaultTip)

	runtimeDispatchInfo := primitives.RuntimeDispatchInfo{
		Weight:     dispatchInfo.Weight,
		Class:      dispatchInfo.Class,
		PartialFee: partialFee,
	}

	return utils.BytesToOffsetAndSize(runtimeDispatchInfo.Bytes())
}

// QueryCallFeeDetails queries the detailed fee of a dispatch call.
// It takes two arguments:
// - dataPtr: Pointer to the data in the Wasm memory.
// - dataLen: Length of the data.
// which represent the SCALE-encoded dispatch call and its length.
// Returns a pointer-size of the SCALE-encoded detailed fee.
// [Specification](https://spec.polkadot.network/chap-runtime-api#sect-rte-transactionpaymentcallapi-query-call-fee-details)
func QueryCallFeeDetails(dataPtr int32, dataLen int32) int64 {
	b := utils.ToWasmMemorySlice(dataPtr, dataLen)
	buffer := bytes.NewBuffer(b)

	call := types.DecodeCall(buffer)
	length := sc.DecodeU32(buffer)

	dispatchInfo := primitives.GetDispatchInfo(call)
	feeDetails := computeFeeDetails(length, dispatchInfo, DefaultTip)

	return utils.BytesToOffsetAndSize(feeDetails.Bytes())
}

func computeFee(len sc.U32, info primitives.DispatchInfo, tip primitives.Balance) primitives.Balance {
	return computeFeeDetails(len, info, tip).FinalFee()
}

func computeFeeDetails(len sc.U32, info primitives.DispatchInfo, tip primitives.Balance) primitives.FeeDetails {
	return computeFeeRaw(len, info.Weight, tip, info.PaysFee, info.Class)
}

func computeActualFee(len sc.U32, info primitives.DispatchInfo, postInfo primitives.PostDispatchInfo, tip primitives.Balance) primitives.Balance {
	return computeActualFeeDetails(len, info, postInfo, tip).FinalFee()
}

func computeActualFeeDetails(len sc.U32, info primitives.DispatchInfo, postInfo primitives.PostDispatchInfo, tip primitives.Balance) primitives.FeeDetails {
	return computeFeeRaw(len, postInfo.CalcActualWeight(&info), tip, postInfo.Pays(&info), info.Class)
}

func computeFeeRaw(len sc.U32, weight primitives.Weight, tip primitives.Balance, paysFee primitives.Pays, class primitives.DispatchClass) primitives.FeeDetails {
	if paysFee[0] == primitives.PaysYes { // TODO: type safety
		unadjustedWeightFee := weightToFee(weight)
		multiplier := storageNextFeeMultiplier()

		fixedU128Div := big.NewInt(1_000_000_000_000_000_000)
		bnAdjustedWeightFee := new(big.Int).Mul(multiplier.ToBigInt(), unadjustedWeightFee.ToBigInt())
		adjustedWeightFee := sc.NewU128FromBigInt(new(big.Int).Div(bnAdjustedWeightFee, fixedU128Div)) // TODO: Create FixedU128 type

		lenFee := lengthToFee(len)
		baseFee := weightToFee(system.DefaultBlockWeights().Get(class).BaseExtrinsic)

		inclusionFee := sc.NewOption[primitives.InclusionFee](primitives.NewInclusionFee(baseFee, lenFee, adjustedWeightFee))

		return primitives.FeeDetails{
			InclusionFee: inclusionFee,
			Tip:          tip,
		}
	}

	return primitives.FeeDetails{
		InclusionFee: sc.NewOption[primitives.InclusionFee](nil),
		Tip:          tip,
	}
}

func lengthToFee(length sc.U32) primitives.Balance {
	return constants.LengthToFee.WeightToFee(primitives.WeightFromParts(sc.U64(length), 0))
}

func weightToFee(weight primitives.Weight) primitives.Balance {
	cappedWeight := weight.Min(system.DefaultBlockWeights().MaxBlock)

	return constants.WeightToFee.WeightToFee(cappedWeight)
}

func storageNextFeeMultiplier() sc.U128 {
	// Storage value is FixedU128, which is different from U128.
	// It implements a decimal fixed point number, which is `1 / VALUE`
	// Example: FixedU128, VALUE is 1_000_000_000_000_000_000.
	// FixedU64, VALUE is 1_000_000_000.
	txPaymentHash := hashing.Twox128(constants.KeyTransactionPayment)
	nextFeeMultiplierHash := hashing.Twox128(constants.KeyNextFeeMultiplier)
	key := append(txPaymentHash, nextFeeMultiplierHash...)

	return storage.GetDecodeOnEmpty(key, sc.DecodeU128, DefaultMultiplierValue)
}
