package types

import (
	"bytes"
	"math/big"

	sc "github.com/LimeChain/goscale"
)

type RefCount = sc.U32

type AccountInfo struct {
	Nonce       AccountIndex
	Consumers   RefCount
	Providers   RefCount
	Sufficients RefCount
	Data        AccountData
}

func (ai AccountInfo) Encode(buffer *bytes.Buffer) {
	ai.Nonce.Encode(buffer)
	ai.Consumers.Encode(buffer)
	ai.Providers.Encode(buffer)
	ai.Sufficients.Encode(buffer)
	ai.Data.Encode(buffer)
}

func (ai AccountInfo) Bytes() []byte {
	return sc.EncodedBytes(ai)
}

func DecodeAccountInfo(buffer *bytes.Buffer) AccountInfo {
	return AccountInfo{
		Nonce:       sc.DecodeU32(buffer),
		Consumers:   sc.DecodeU32(buffer),
		Providers:   sc.DecodeU32(buffer),
		Sufficients: sc.DecodeU32(buffer),
		Data:        DecodeAccountData(buffer),
	}
}

func (ai AccountInfo) Frozen(reasons Reasons) *big.Int {
	switch reasons {
	case ReasonsAll:
		if ai.Data.MiscFrozen.ToBigInt().Cmp(ai.Data.FeeFrozen.ToBigInt()) > 0 {
			return ai.Data.MiscFrozen.ToBigInt()
		}
		return ai.Data.FeeFrozen.ToBigInt()
	case ReasonsMisc:
		return big.NewInt(0).Set(ai.Data.MiscFrozen.ToBigInt())
	case ReasonsFee:
		return big.NewInt(0).Set(ai.Data.MiscFrozen.ToBigInt())
	}

	return big.NewInt(0)
}
