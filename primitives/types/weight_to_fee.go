package types

type WeightToFee interface {
	WeightToFee(weight Weight) Balance
}
