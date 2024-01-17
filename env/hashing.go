//go:build !nonwasmenv

package env

/*
	Hashing: Interface that provides functions for hashing with different algorithms.
*/

//go:wasm-module env
//go:export ext_hashing_blake2_128_version_1
func ExtHashingBlake2128Version1(data int64) int32

//go:wasm-module env
//go:export ext_hashing_blake2_256_version_1
func ExtHashingBlake2256Version1(data int64) int32

//go:wasm-module env
//go:export ext_hashing_keccak_256_version_1
func ExtHashingKeccak256Version1(data int64) int32

//go:wasm-module env
//go:export ext_hashing_twox_128_version_1
func ExtHashingTwox128Version1(data int64) int32

//go:wasm-module env
//go:export ext_hashing_twox_64_version_1
func ExtHashingTwox64Version1(data int64) int32
