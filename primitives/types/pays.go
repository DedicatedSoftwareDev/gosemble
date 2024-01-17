package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/log"
)

const (
	// PaysYes Transactor will pay related fees.
	PaysYes sc.U8 = iota

	// PaysNo Transactor will NOT pay related fees.
	PaysNo
)

type Pays = sc.VaryingData

func NewPaysYes() Pays {
	return sc.NewVaryingData(PaysYes)
}

func NewPaysNo() Pays {
	return sc.NewVaryingData(PaysNo)
}

func DecodePays(buffer *bytes.Buffer) Pays {
	b := sc.DecodeU8(buffer)

	switch b {
	case PaysYes:
		return NewPaysYes()
	case PaysNo:
		return NewPaysNo()
	default:
		log.Critical("invalid Pays type")
	}

	panic("unreachable")
}
