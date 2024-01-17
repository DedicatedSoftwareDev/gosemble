//go:build !nonwasmenv

package trie

import (
	"github.com/LimeChain/gosemble/env"
	"github.com/LimeChain/gosemble/utils"
)

func Blake2256OrderedRoot(key []byte, version int32) []byte {
	keyOffsetSize := utils.BytesToOffsetAndSize(key)
	r := env.ExtTrieBlake2256OrderedRootVersion2(keyOffsetSize, version)
	return utils.ToWasmMemorySlice(r, 32)
}
