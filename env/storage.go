//go:build !nonwasmenv

package env

/*
	Storage: Interface for manipulating the storage from within the runtime.
*/

//go:wasm-module env
//go:export ext_storage_append_version_1
func ExtStorageAppendVersion1(key int64, value int64)

//go:wasm-module env
//go:export ext_storage_clear_version_1
func ExtStorageClearVersion1(key_data int64)

//go:wasm-module env
//go:export ext_storage_clear_prefix_version_2
func ExtStorageClearPrefixVersion2(prefix int64, limit int64) int64

//go:wasm-module env
//go:export ext_storage_exists_version_1
func ExtStorageExistsVersion1(key int64) int32

//go:wasm-module env
//go:export ext_storage_get_version_1
func ExtStorageGetVersion1(key int64) int64

//go:wasm-module env
//go:export ext_storage_next_key_version_1
func ExtStorageNextKeyVersion1(key int64) int64

//go:wasm-module env
//go:export ext_storage_read_version_1
func ExtStorageReadVersion1(key int64, value_out int64, offset int32) int64

//go:wasm-module env
//go:export ext_storage_root_version_2
func ExtStorageRootVersion2(key int32) int64

//go:wasm-module env
//go:export ext_storage_set_version_1
func ExtStorageSetVersion1(key int64, value int64)

//go:wasm-module env
//go:export ext_storage_start_transaction_version_1
func ExtStorageStartTransactionVersion1()

//go:wasm-module env
//go:export ext_storage_commit_transaction_version_1
func ExtStorageCommitTransactionVersion()

//go:wasm-module env
//go:export ext_storage_rollback_transaction_version_1
func ExtStorageRollbackTransactionVersion1()
