package constants

import "math/big"

const (
	MilliCents        = Cents / 1000
	Cents             = Dollar / 100
	Dollar            = Units
	Units      uint64 = 10_000_000_000
)

var (
	Zero = big.NewInt(0)
)
