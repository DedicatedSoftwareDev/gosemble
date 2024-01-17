package timestamp

import (
	"github.com/LimeChain/gosemble/constants"
	"github.com/LimeChain/gosemble/primitives/hashing"
	"github.com/LimeChain/gosemble/primitives/log"
	"github.com/LimeChain/gosemble/primitives/storage"
)

func OnFinalize() {
	timestampHash := hashing.Twox128(constants.KeyTimestamp)
	didUpdateHash := hashing.Twox128(constants.KeyDidUpdate)

	didUpdate := storage.Get(append(timestampHash, didUpdateHash...))

	if didUpdate.HasValue {
		storage.Clear(append(timestampHash, didUpdateHash...))
	} else {
		log.Critical("Timestamp must be updated once in the block")
	}
}
