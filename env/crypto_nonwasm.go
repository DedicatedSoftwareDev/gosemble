//go:build nonwasmenv

package env

/*
	Crypto: Interfaces for working with crypto related types from within the runtime.
*/

func ExtCryptoEd25519GenerateVersion1(key_type_id int32, seed int64) int32 {
	panic("not implemented")
}

func ExtCryptoEd25519VerifyVersion1(sig int32, msg int64, key int32) int32 {
	panic("not implemented")
}

func ExtCryptoFinishBatchVerifyVersion1() int32 {
	panic("not implemented")
}

func ExtCryptoSecp256k1EcdsaRecoverVersion2(sig int32, msg int32) int64 {
	panic("not implemented")
}

func ExtCryptoSecp256k1EcdsaRecoverCompressedVersion2(sig int32, msg int32) int64 {
	panic("not implemented")
}

func ExtCryptoSr25519GenerateVersion1(key_type_id int32, seed int64) int32 {
	panic("not implemented")
}

func ExtCryptoSr25519PublicKeysVersion1(key_type_id int32) int64 {
	panic("not implemented")
}

func ExtCryptoSr25519SignVersion1(key_type_id int32, key int32, msg int64) int64 {
	panic("not implemented")
}

func ExtCryptoSr25519VerifyVersion2(sig int32, msg int64, key int32) int32 {
	panic("not implemented")
}

func ExtCryptoStartBatchVerifyVersion1() {
	panic("not implemented")
}
