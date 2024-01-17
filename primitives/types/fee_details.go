package types

import (
	"bytes"
	"math/big"

	sc "github.com/LimeChain/goscale"
)

type FeeDetails struct {
	InclusionFee sc.Option[InclusionFee]

	Tip Balance // not serializable
}

func (fd FeeDetails) Encode(buffer *bytes.Buffer) {
	fd.InclusionFee.Encode(buffer)
}

func (fd FeeDetails) Bytes() []byte {
	return sc.EncodedBytes(fd)
}

func DecodeFeeDetails(buffer *bytes.Buffer) FeeDetails {
	return FeeDetails{
		InclusionFee: sc.DecodeOptionWith(buffer, DecodeInclusionFee),
	}
}

func (fd FeeDetails) FinalFee() Balance {
	sum := fd.Tip
	if fd.InclusionFee.HasValue {
		inclusionFee := fd.InclusionFee.Value.InclusionFee().ToBigInt()
		total := new(big.Int).Add(inclusionFee, fd.Tip.ToBigInt())

		sum = sc.NewU128FromBigInt(total)
	}

	return sum
}
