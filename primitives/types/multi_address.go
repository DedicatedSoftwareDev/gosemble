package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/log"
)

// AccountId It's an account ID (pubkey).
type AccountId struct {
	Address32 // TODO: Varies depending on Signature (32 for ed25519 and sr25519, 33 for ecdsa)
}

func DecodeAccountId(buffer *bytes.Buffer) AccountId {
	return AccountId{DecodeAddress32(buffer)} // TODO: length 32 or 33 depending on algorithm
}

// AccountIndex It's an account index.
type AccountIndex = sc.U32

// AccountRaw It's some arbitrary raw bytes.
type AccountRaw struct {
	sc.Sequence[sc.U8]
}

func (a AccountRaw) Encode(buffer *bytes.Buffer) {
	a.Sequence.Encode(buffer)
}

func DecodeAccountRaw(buffer *bytes.Buffer) AccountRaw {
	return AccountRaw{sc.DecodeSequence[sc.U8](buffer)}
}

// Address32 It's a 32 byte representation.
type Address32 struct {
	sc.FixedSequence[sc.U8] // size 32
}

func NewAddress32(values ...sc.U8) Address32 {
	if len(values) != 32 {
		log.Critical("Address32 should be of size 32")
	}
	return Address32{sc.NewFixedSequence(32, values...)}
}

func DecodeAddress32(buffer *bytes.Buffer) Address32 {
	return Address32{sc.DecodeFixedSequence[sc.U8](32, buffer)}
}

// Address20 It's a 20 byte representation.
type Address20 struct {
	sc.FixedSequence[sc.U8] // size 20
}

func NewAddress20(values ...sc.U8) Address20 {
	if len(values) != 20 {
		log.Critical("Address20 should be of size 20")
	}
	return Address20{sc.NewFixedSequence(20, values...)}
}

func DecodeAddress20(buffer *bytes.Buffer) Address20 {
	return Address20{sc.DecodeFixedSequence[sc.U8](20, buffer)}
}

const (
	MultiAddressId sc.U8 = iota
	MultiAddressIndex
	MultiAddressRaw
	MultiAddress32
	MultiAddress20
)

type MultiAddress struct {
	sc.VaryingData
}

func NewMultiAddressId(id AccountId) MultiAddress {
	return MultiAddress{sc.NewVaryingData(MultiAddressId, id)}
}

func NewMultiAddressIndex(index AccountIndex) MultiAddress {
	return MultiAddress{sc.NewVaryingData(MultiAddressIndex, sc.ToCompact(index))}
}

func NewMultiAddressRaw(accountRaw AccountRaw) MultiAddress {
	return MultiAddress{sc.NewVaryingData(MultiAddressRaw, accountRaw)}
}

func NewMultiAddress32(address Address32) MultiAddress {
	return MultiAddress{sc.NewVaryingData(MultiAddress32, address)}
}

func NewMultiAddress20(address Address20) MultiAddress {
	return MultiAddress{sc.NewVaryingData(MultiAddress20, address)}
}

func DecodeMultiAddress(buffer *bytes.Buffer) MultiAddress {
	b := sc.DecodeU8(buffer)

	switch b {
	case MultiAddressId:
		return NewMultiAddressId(DecodeAccountId(buffer))
	case MultiAddressIndex:
		compact := sc.DecodeCompact(buffer)
		index := sc.U32(compact.ToBigInt().Int64())
		return NewMultiAddressIndex(index)
	case MultiAddressRaw:
		return NewMultiAddressRaw(DecodeAccountRaw(buffer))
	case MultiAddress32:
		return NewMultiAddress32(DecodeAddress32(buffer))
	case MultiAddress20:
		return NewMultiAddress20(DecodeAddress20(buffer))
	default:
		log.Critical("invalid MultiAddress type in Decode")
	}

	panic("unreachable")
}

func (a MultiAddress) IsAccountId() sc.Bool {
	switch a.VaryingData[0] {
	case MultiAddressId:
		return true
	default:
		return false
	}
}

func (a MultiAddress) AsAccountId() AccountId {
	if a.IsAccountId() {
		return a.VaryingData[1].(AccountId)
	} else {
		log.Critical("not an AccountId type")
	}

	panic("unreachable")
}

func (a MultiAddress) IsAccountIndex() sc.Bool {
	switch a.VaryingData[0] {
	case MultiAddressIndex:
		return true
	default:
		return false
	}
}

func (a MultiAddress) AsAccountIndex() AccountIndex {
	if a.IsAccountIndex() {
		compact := a.VaryingData[1].(sc.Compact).ToBigInt()

		return sc.U32(compact.Int64())
	} else {
		log.Critical("not an AccountIndex type")
	}

	panic("unreachable")
}

func (a MultiAddress) IsRaw() sc.Bool {
	switch a.VaryingData[0] {
	case MultiAddressRaw:
		return true
	default:
		return false
	}
}

func (a MultiAddress) AsRaw() AccountRaw {
	if a.IsRaw() {
		return a.VaryingData[1].(AccountRaw)
	} else {
		log.Critical("not an AccountRaw type")
	}

	panic("unreachable")
}

func (a MultiAddress) IsAddress32() sc.Bool {
	switch a.VaryingData[0] {
	case MultiAddress32:
		return true
	default:
		return false
	}
}

func (a MultiAddress) AsAddress32() Address32 {
	if a.IsAddress32() {
		return a.VaryingData[1].(Address32)
	} else {
		log.Critical("not an Address32 type")
	}

	panic("unreachable")
}

func (a MultiAddress) IsAddress20() sc.Bool {
	switch a.VaryingData[0] {
	case MultiAddress20:
		return true
	default:
		return false
	}
}

func (a MultiAddress) AsAddress20() Address20 {
	if a.IsAddress20() {
		return a.VaryingData[1].(Address20)
	} else {
		log.Critical("not an Address20 type")
	}

	panic("unreachable")
}
