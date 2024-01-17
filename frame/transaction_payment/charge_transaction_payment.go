package transaction_payment

import (
	"math/big"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants"
	"github.com/LimeChain/gosemble/constants/transaction_payment"
	"github.com/LimeChain/gosemble/frame/balances/dispatchables"
	"github.com/LimeChain/gosemble/frame/system"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

type ChargeTransactionPayment primitives.Balance

func (ctp ChargeTransactionPayment) AdditionalSigned() (ok sc.Empty, err primitives.TransactionValidityError) {
	return sc.Empty{}, nil
}

func (ctp ChargeTransactionPayment) Validate(who *primitives.Address32, call *primitives.Call, info *primitives.DispatchInfo, length sc.Compact) (primitives.ValidTransaction, primitives.TransactionValidityError) {
	finalFee, _, err := ctp.withdrawFee(who, call, info, length)
	if err != nil {
		return primitives.ValidTransaction{}, err
	}

	tip := primitives.Balance(ctp)
	validTransaction := primitives.DefaultValidTransaction()
	validTransaction.Priority = ctp.getPriority(info, length, tip, finalFee)

	return validTransaction, nil
}

func (ctp ChargeTransactionPayment) PreDispatch(who *primitives.Address32, call *primitives.Call, info *primitives.DispatchInfo, length sc.Compact) (ok primitives.Pre, err primitives.TransactionValidityError) {
	_, imbalance, err := ctp.withdrawFee(who, call, info, length)
	return primitives.Pre{
		Tip:       primitives.Balance(ctp),
		Who:       *who,
		Imbalance: imbalance,
	}, err
}

func (ctp ChargeTransactionPayment) PostDispatch(pre sc.Option[primitives.Pre], info *primitives.DispatchInfo, postInfo *primitives.PostDispatchInfo, length sc.Compact, result *primitives.DispatchResult) (primitives.Pre, primitives.TransactionValidityError) {
	if pre.HasValue {
		preValue := pre.Value
		actualFee := computeActualFee(sc.U32(length.ToBigInt().Uint64()), *info, *postInfo, preValue.Tip)
		err := correctAndDepositFee(&preValue.Who, actualFee, preValue.Tip, preValue.Imbalance)
		if err != nil {
			return primitives.Pre{}, err
		}

		system.DepositEvent(NewEventTransactionFeePaid(preValue.Who.FixedSequence, actualFee, preValue.Tip))
	}
	return primitives.Pre{}, nil
}

func (ctp ChargeTransactionPayment) getPriority(info *primitives.DispatchInfo, len sc.Compact, tip primitives.Balance, finalFee primitives.Balance) primitives.TransactionPriority {
	maxBlockWeight := system.DefaultBlockWeights().MaxBlock.RefTime
	maxDefaultBlockLength := system.DefaultBlockLength().Max
	maxBlockLength := sc.U64(*maxDefaultBlockLength.Get(info.Class))

	infoWeight := info.Weight.RefTime

	// info_weight.clamp(1, max_block_weight);
	boundedWeight := infoWeight
	if boundedWeight < 1 {
		boundedWeight = 1
	} else if boundedWeight > maxBlockWeight {
		boundedWeight = maxBlockWeight
	}

	// (len as u64).clamp(1, max_block_length);
	boundedLength := sc.U64(len.ToBigInt().Uint64())
	if boundedLength < 1 {
		boundedLength = 1
	} else if boundedLength > maxBlockLength {
		boundedLength = maxBlockLength
	}

	maxTxPerBlockWeight := maxBlockWeight / boundedWeight
	maxTxPerBlockLength := maxBlockLength / boundedLength

	maxTxPerBlock := maxTxPerBlockWeight
	if maxTxPerBlockWeight > maxTxPerBlockLength {
		maxTxPerBlock = maxTxPerBlockLength
	}

	bnTip := new(big.Int).Add(tip.ToBigInt(), big.NewInt(1))

	scaledTip := new(big.Int).Mul(bnTip, new(big.Int).SetUint64(uint64(maxTxPerBlock)))

	if info.Class.Is(primitives.DispatchClassNormal) {
		return sc.U64(scaledTip.Uint64())
	} else if info.Class.Is(primitives.DispatchClassMandatory) {
		return sc.U64(scaledTip.Uint64())
	} else if info.Class.Is(primitives.DispatchClassOperational) {
		feeMultiplier := transaction_payment.OperationalFeeMultiplier
		virtualTip := new(big.Int).Mul(finalFee.ToBigInt(), big.NewInt(int64(feeMultiplier)))
		scaledVirtualTip := new(big.Int).Mul(virtualTip, new(big.Int).SetUint64(uint64(maxTxPerBlock)))

		sum := new(big.Int).Add(scaledTip, scaledVirtualTip)

		return sc.U64(sum.Uint64())
	}

	return 0
}

func (ctp ChargeTransactionPayment) withdrawFee(who *primitives.Address32, _call *primitives.Call, info *primitives.DispatchInfo, length sc.Compact) (primitives.Balance, sc.Option[primitives.Balance], primitives.TransactionValidityError) {
	tip := primitives.Balance(ctp)
	fee := computeFee(sc.U32(length.ToBigInt().Uint64()), *info, tip)

	imbalance, err := withdrawFee(who, _call, info, fee, sc.NewU128FromBigInt(tip.ToBigInt()))
	if err != nil {
		return primitives.Balance{}, sc.NewOption[primitives.Balance](nil), err
	}

	return fee, imbalance, nil
}

func withdrawFee(who *primitives.Address32, _call *primitives.Call, _info *primitives.DispatchInfo, fee primitives.Balance, tip primitives.Balance) (sc.Option[primitives.Balance], primitives.TransactionValidityError) {
	if fee.ToBigInt().Cmp(constants.Zero) == 0 {
		return sc.NewOption[primitives.Balance](nil), nil
	}

	withdrawReasons := primitives.WithdrawReasonsTransactionPayment
	if tip.ToBigInt().Cmp(constants.Zero) == 0 {
		withdrawReasons = primitives.WithdrawReasonsTransactionPayment
	} else {
		withdrawReasons = primitives.WithdrawReasonsTransactionPayment | primitives.WithdrawReasonsTip
	}

	imbalance, err := dispatchables.Withdraw(*who, fee, sc.U8(withdrawReasons), primitives.ExistenceRequirementKeepAlive)
	if err != nil {
		return sc.NewOption[primitives.Balance](nil), primitives.NewTransactionValidityError(primitives.NewInvalidTransactionPayment())
	}

	return sc.NewOption[primitives.Balance](imbalance), nil
}

func correctAndDepositFee(who *primitives.Address32, correctedFee primitives.Balance, tip primitives.Balance, alreadyWithdrawn sc.Option[primitives.Balance]) primitives.TransactionValidityError {
	if alreadyWithdrawn.HasValue {
		alreadyPaidNegativeImbalance := alreadyWithdrawn.Value
		refundAmount := new(big.Int).Sub(alreadyPaidNegativeImbalance.ToBigInt(), correctedFee.ToBigInt())

		refundPositiveImbalance, err := dispatchables.DepositIntoExisting(*who, sc.NewU128FromBigInt(refundAmount))
		if err != nil {
			return primitives.NewTransactionValidityError(primitives.NewInvalidTransactionPayment())
		}

		comparison := alreadyPaidNegativeImbalance.ToBigInt().Cmp(refundPositiveImbalance.ToBigInt())
		if comparison < 0 {
			return primitives.NewTransactionValidityError(primitives.NewInvalidTransactionPayment())
		}
	}
	return nil
}
