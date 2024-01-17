//go:build !nonwasmenv

package env

/*
	Trie: Interface that provides trie related functionality
*/

//go:wasm-module env
//go:export ext_trie_blake2_256_ordered_root_version_2
func ExtTrieBlake2256OrderedRootVersion2(input int64, version int32) int32
