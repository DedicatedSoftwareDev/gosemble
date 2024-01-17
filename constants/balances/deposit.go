package balances

import (
	"math/big"

	"github.com/LimeChain/gosemble/constants"
)

const (
	MaxLocks    = 50
	MaxReserves = 50
)

var (
	existentialDeposit = 1 * constants.Dollar
	ExistentialDeposit = big.NewInt(0).SetUint64(existentialDeposit)
)
