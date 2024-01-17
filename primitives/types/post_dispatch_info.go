package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
)

// PostDispatchInfo Weight information that is only available post dispatch.
// NOTE: This can only be used to reduce the weight or fee, not increase it.
type PostDispatchInfo struct {
	// Actual weight consumed by a call or `None` which stands for the worst case static weight.
	ActualWeight sc.Option[Weight]

	// Whether this transaction should pay fees when all is said and done.
	PaysFee sc.U8
}

func (pdi PostDispatchInfo) Encode(buffer *bytes.Buffer) {
	pdi.ActualWeight.Encode(buffer)
	pdi.PaysFee.Encode(buffer)
}

func DecodePostDispatchInfo(buffer *bytes.Buffer) PostDispatchInfo {
	pdi := PostDispatchInfo{}
	pdi.ActualWeight = sc.DecodeOptionWith(buffer, DecodeWeight)
	pdi.PaysFee = sc.DecodeU8(buffer)
	return pdi
}

func (pdi PostDispatchInfo) Bytes() []byte {
	return sc.EncodedBytes(pdi)
}

// CalcUnspent Calculate how much (if any) weight was not used by the `Dispatchable`.
func (pdi PostDispatchInfo) CalcUnspent(info *DispatchInfo) Weight {
	return info.Weight.Sub(pdi.CalcActualWeight(info))
}

// CalcActualWeight Calculate how much weight was actually spent by the `Dispatchable`.
func (pdi PostDispatchInfo) CalcActualWeight(info *DispatchInfo) Weight {
	if pdi.ActualWeight.HasValue {
		actualWeight := pdi.ActualWeight.Value
		return actualWeight.Min(info.Weight)
	} else {
		return info.Weight
	}
}

// Pays Determine if user should actually pay fees at the end of the dispatch.
func (pdi PostDispatchInfo) Pays(info *DispatchInfo) Pays {
	// If they originally were not paying fees, or the post dispatch info
	// says they should not pay fees, then they don't pay fees.
	// This is because the pre dispatch information must contain the
	// worst case for weight and fees paid.

	if info.PaysFee[0] == PaysNo || pdi.PaysFee == PaysNo {
		return NewPaysNo()
	} else {
		// Otherwise they pay.
		return NewPaysYes()
	}
}
