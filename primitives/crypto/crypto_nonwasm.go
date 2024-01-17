//go:build nonwasmenv

package crypto

import (
	"crypto/ed25519"

	sc "github.com/LimeChain/goscale"
)

func ExtCryptoEd25519GenerateVersion1(keyTypeId []byte, seed []byte) []byte {
	panic("not implemented")
}

func ExtCryptoEd25519VerifyVersion1(signature []byte, message []byte, pubKey []byte) sc.Bool {
	return sc.Bool(ed25519.Verify(pubKey, message, signature))
}

func ExtCryptoSr25519GenerateVersion1(keyTypeId []byte, seed []byte) []byte {
	panic("not implemented")
}

func ExtCryptoSr25519VerifyVersion2(signature []byte, message []byte, pubKey []byte) sc.Bool {
	panic("not implemented")
}

func ExtCryptoStartBatchVerify() {
	panic("not implemented")
}

func ExtCryptoFinishBatchVerify() int32 {
	panic("not implemented")
}
