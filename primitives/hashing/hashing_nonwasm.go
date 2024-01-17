//go:build nonwasmenv

package hashing

import (
	"github.com/ChainSafe/gossamer/lib/common"
)

func Twox128(value []byte) []byte {
	h, _ := common.Twox128Hash(value)
	return h[:]
}

func Twox64(value []byte) []byte {
	h, _ := common.Twox64(value)
	return h[:]
}

func Blake128(value []byte) []byte {
	h, _ := common.Blake2b128(value)
	return h[:]
}

func Blake256(value []byte) []byte {
	h, _ := common.Blake2bHash(value)
	return h[:]
}
