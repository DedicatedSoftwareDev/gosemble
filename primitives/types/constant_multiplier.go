package types

import (
	"math/big"

	sc "github.com/LimeChain/goscale"
)

type ConstantMultiplier struct {
	Multiplier Balance
}

func NewConstantMultiplier(multiplier Balance) ConstantMultiplier {
	return ConstantMultiplier{
		Multiplier: multiplier,
	}
}

func (cm ConstantMultiplier) WeightToFee(weight Weight) Balance {
	bnRefTime := new(big.Int).SetUint64(uint64(weight.RefTime))

	res := new(big.Int).Mul(bnRefTime, cm.Multiplier.ToBigInt())
	return sc.NewU128FromBigInt(res)
}
