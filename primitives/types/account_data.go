package types

import (
	"bytes"
	"math/big"

	sc "github.com/LimeChain/goscale"
)

type Balance = sc.U128

type AccountData struct {
	Free       Balance
	Reserved   Balance
	MiscFrozen Balance
	FeeFrozen  Balance
}

func (ad AccountData) Encode(buffer *bytes.Buffer) {
	ad.Free.Encode(buffer)
	ad.Reserved.Encode(buffer)
	ad.MiscFrozen.Encode(buffer)
	ad.FeeFrozen.Encode(buffer)
}

func (ad AccountData) Bytes() []byte {
	return sc.EncodedBytes(ad)
}

func DecodeAccountData(buffer *bytes.Buffer) AccountData {
	return AccountData{
		Free:       sc.DecodeU128(buffer),
		Reserved:   sc.DecodeU128(buffer),
		MiscFrozen: sc.DecodeU128(buffer),
		FeeFrozen:  sc.DecodeU128(buffer),
	}
}

func (ad AccountData) Total() *big.Int {
	return new(big.Int).Add(ad.Free.ToBigInt(), ad.Reserved.ToBigInt())
}
