package types

const (
	WithdrawReasonsTransactionPayment = 1 << iota
	WithdrawReasonsTransfer
	WithdrawReasonsReserve
	WithdrawReasonsFee
	WithdrawReasonsTip
)
