//go:build nonwasmenv

package env

/*
	Storage: Interface for manipulating the storage from within the runtime.
*/

func ExtStorageAppendVersion1(key int64, value int64) int64 {
	panic("not implemented")
}

func ExtStorageClearVersion1(key_data int64) {
	panic("not implemented")
}

func ExtStorageClearPrefixVersion2(prefix int64, limit int64) int64 {
	panic("not implemented")
}

func ExtStorageCommitTransactionVersion() {
	panic("not implemented")
}

func ExtStorageExistsVersion1(key int64) int32 {
	panic("not implemented")
}

func ExtStorageGetVersion1(key int64) int64 {
	panic("not implemented")
}

func ExtStorageNextKeyVersion1(key int64) int64 {
	panic("not implemented")
}

func ExtStorageReadVersion1(key int64, value_out int64, offset int32) int64 {
	panic("not implemented")
}

func ExtStorageRollbackTransactionVersion1() {
	panic("not implemented")
}

func ExtStorageRootVersion2(key int32) int64 {
	panic("not implemented")
}

func ExtStorageSetVersion1(key int64, value int64) {
	panic("not implemented")
}

func ExtStorageStartTransactionVersion1() {
	panic("not implemented")
}
