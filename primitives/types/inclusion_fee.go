package types

import (
	"bytes"
	"math/big"

	sc "github.com/LimeChain/goscale"
)

type InclusionFee struct {
	BaseFee           Balance
	LenFee            Balance
	AdjustedWeightFee Balance
}

func NewInclusionFee(baseFee, lenFee, adjustedWeightFee Balance) InclusionFee {
	return InclusionFee{
		baseFee,
		lenFee,
		adjustedWeightFee,
	}
}

func (i InclusionFee) Encode(buffer *bytes.Buffer) {
	i.BaseFee.Encode(buffer)
	i.LenFee.Encode(buffer)
	i.AdjustedWeightFee.Encode(buffer)
}

func (i InclusionFee) Bytes() []byte {
	return sc.EncodedBytes(i)
}

func DecodeInclusionFee(buffer *bytes.Buffer) InclusionFee {
	return InclusionFee{
		BaseFee:           sc.DecodeU128(buffer),
		LenFee:            sc.DecodeU128(buffer),
		AdjustedWeightFee: sc.DecodeU128(buffer),
	}
}

func (i InclusionFee) InclusionFee() Balance {
	sum := new(big.Int).Add(i.BaseFee.ToBigInt(), i.LenFee.ToBigInt())

	sum = sum.Add(sum, i.AdjustedWeightFee.ToBigInt())

	return sc.NewU128FromBigInt(sum)
}
