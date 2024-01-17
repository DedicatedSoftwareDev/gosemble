package types

import (
	"bytes"

	"github.com/LimeChain/gosemble/primitives/log"

	sc "github.com/LimeChain/goscale"
)

const (
	BalanceStatusFree sc.U8 = iota
	BalanceStatusReserved
)

type BalanceStatus = sc.U8

func DecodeBalanceStatus(buffer *bytes.Buffer) sc.U8 {
	value := sc.DecodeU8(buffer)
	switch value {
	case BalanceStatusFree, BalanceStatusReserved:
		return value
	default:
		log.Critical("invalid balance status type")
	}

	panic("unreachable")
}
