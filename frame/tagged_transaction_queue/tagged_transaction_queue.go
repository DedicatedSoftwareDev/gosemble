package tagged_transaction_queue

import (
	"bytes"

	"github.com/LimeChain/gosemble/execution/types"
	"github.com/LimeChain/gosemble/frame/executive"
	primitives "github.com/LimeChain/gosemble/primitives/types"
	"github.com/LimeChain/gosemble/utils"
)

type TaggedTransactionQueue interface {
	ValidateTransaction(dataPtr int32, dataLen int32) int64
}

// ValidateTransaction validates an extrinsic at a given block.
// It takes two arguments:
// - dataPtr: Pointer to the data in the Wasm memory.
// - dataLen: Length of the data.
// which represent the SCALE-encoded tx source, extrinsic and block hash.
// Returns a pointer-size of the SCALE-encoded result whether the extrinsic is valid.
// [Specification](https://spec.polkadot.network/#sect-rte-validate-transaction)
func ValidateTransaction(dataPtr int32, dataLen int32) int64 {
	data := utils.ToWasmMemorySlice(dataPtr, dataLen)
	buffer := bytes.NewBuffer(data)

	txSource := primitives.DecodeTransactionSource(buffer)
	tx := types.DecodeUncheckedExtrinsic(buffer)
	blockHash := primitives.DecodeBlake2bHash(buffer)

	ok, err := executive.ValidateTransaction(txSource, tx, blockHash)

	var res primitives.TransactionValidityResult
	if err != nil {
		res = primitives.NewTransactionValidityResult(err)
	} else {
		res = primitives.NewTransactionValidityResult(ok)
	}

	return utils.BytesToOffsetAndSize(res.Bytes())
}
