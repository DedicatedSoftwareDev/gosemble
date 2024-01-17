package errors

import sc "github.com/LimeChain/goscale"

// Balances module errors.
const (
	ErrorVestingBalance sc.U8 = iota
	ErrorLiquidityRestrictions
	ErrorInsufficientBalance
	ErrorExistentialDeposit
	ErrorKeepAlive
	ErrorExistingVestingSchedule
	ErrorDeadAccount
	ErrorTooManyReserves
)
