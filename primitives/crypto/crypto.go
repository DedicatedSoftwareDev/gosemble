//go:build !nonwasmenv

package crypto

import (
	"github.com/LimeChain/gosemble/env"
	"github.com/LimeChain/gosemble/utils"
)

func ExtCryptoEd25519GenerateVersion1(keyTypeId []byte, seed []byte) []byte {
	r := env.ExtCryptoEd25519GenerateVersion1(utils.Offset32(keyTypeId), utils.BytesToOffsetAndSize(seed))
	return utils.ToWasmMemorySlice(r, 32)
}

func ExtCryptoEd25519VerifyVersion1(signature []byte, message []byte, pubKey []byte) bool {
	return env.ExtCryptoEd25519VerifyVersion1(
		argsSigMsgPubKeyAsWasmMemory(signature, message, pubKey),
	) == 1
}

func ExtCryptoSr25519GenerateVersion1(keyTypeId []byte, seed []byte) []byte {
	r := env.ExtCryptoSr25519GenerateVersion1(utils.Offset32(keyTypeId), utils.BytesToOffsetAndSize(seed))
	return utils.ToWasmMemorySlice(r, 32)
}

func ExtCryptoSr25519VerifyVersion2(signature []byte, message []byte, pubKey []byte) bool {
	return env.ExtCryptoSr25519VerifyVersion2(
		argsSigMsgPubKeyAsWasmMemory(signature, message, pubKey),
	) == 1
}

func ExtCryptoStartBatchVerify() {
	env.ExtCryptoStartBatchVerifyVersion1()
}

func ExtCryptoFinishBatchVerify() int32 {
	return env.ExtCryptoFinishBatchVerifyVersion1()
}

func argsSigMsgPubKeyAsWasmMemory(signature []byte, message []byte, pubKey []byte) (sigOffset int32, msgOffsetSize int64, pubKeyOffset int32) {
	sigOffsetSize := utils.BytesToOffsetAndSize(signature)
	sigOffset, _ = utils.Int64ToOffsetAndSize(sigOffsetSize) // signature: 64-byte

	msgOffsetSize = utils.BytesToOffsetAndSize(message)

	pubKeyOffsetSize := utils.BytesToOffsetAndSize(pubKey)
	pubKeyOffset, _ = utils.Int64ToOffsetAndSize(pubKeyOffsetSize) // public key: 256-bit

	return sigOffset, msgOffsetSize, pubKeyOffset
}
