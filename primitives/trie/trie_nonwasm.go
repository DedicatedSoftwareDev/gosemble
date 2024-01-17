//go:build nonwasmenv

package trie

import (
	"fmt"
	"math/big"

	"github.com/ChainSafe/gossamer/lib/trie"
	"github.com/ChainSafe/gossamer/pkg/scale"
)

func Blake2256OrderedRoot(inherentExt []byte, version int32) []byte {
	var exts [][]byte
	err := scale.Unmarshal(inherentExt, &exts)
	if err != nil {
		panic(err)
	}

	t := trie.NewEmptyTrie()

	for i, value := range exts {
		key, err := scale.Marshal(big.NewInt(int64(i)))
		if err != nil {
			panic(fmt.Sprintf("failed scale encoding value index %d: %s", i, err))
		}

		err = t.Put(key, value)
		if err != nil {
			panic(fmt.Sprintf("failed putting key 0x%x and value 0x%x into trie: %s",
				key, value, err))
		}
	}

	hash, err := t.Hash()
	if err != nil {
		panic(err)
	}

	return hash.ToBytes()
}
