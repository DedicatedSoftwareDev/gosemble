package types

import sc "github.com/LimeChain/goscale"

// IdentityFee implements WeightToFee and maps one unit of weight
// to one unit of fee.
type IdentityFee struct {
}

func (i IdentityFee) WeightToFee(weight Weight) Balance {
	return sc.NewU128FromUint64(uint64(weight.RefTime))
}
