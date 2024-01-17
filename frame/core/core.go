package core

import (
	"bytes"

	"github.com/LimeChain/gosemble/constants"
	"github.com/LimeChain/gosemble/execution/types"
	"github.com/LimeChain/gosemble/frame/executive"
	primitives "github.com/LimeChain/gosemble/primitives/types"
	"github.com/LimeChain/gosemble/utils"
)

type Core interface {
	Version(dataPtr int32, dataLen int32) int64
	ExecuteBlock(dataPtr int32, dataLen int32)
	InitializeBlock(dataPtr int32, dataLen int32)
}

// Version returns a pointer-size SCALE-encoded Runtime version.
// [Specification](https://spec.polkadot.network/#defn-rt-core-version)
func Version() int64 {
	buffer := &bytes.Buffer{}
	constants.RuntimeVersion.Encode(buffer)

	return utils.BytesToOffsetAndSize(buffer.Bytes())
}

// InitializeBlock starts the execution of a particular block.
// It takes two arguments:
// - dataPtr: Pointer to the data in the Wasm memory.
// - dataLen: Length of the data.
// which represent the SCALE-encoded header of the block.
// [Specification](https://spec.polkadot.network/#sect-rte-core-initialize-block)
func InitializeBlock(dataPtr int32, dataLen int32) {
	data := utils.ToWasmMemorySlice(dataPtr, dataLen)
	buffer := bytes.NewBuffer(data)

	header := primitives.DecodeHeader(buffer)
	executive.InitializeBlock(header)
}

// ExecuteBlock executes the provided block.
// It takes two arguments:
// - dataPtr: Pointer to the data in the Wasm memory.
// - dataLen: Length of the data.
// which represent the SCALE-encoded block.
// [Specification](https://spec.polkadot.network/#sect-rte-core-execute-block)
func ExecuteBlock(dataPtr int32, dataLen int32) {
	data := utils.ToWasmMemorySlice(dataPtr, dataLen)
	buffer := bytes.NewBuffer(data)

	block := types.DecodeBlock(buffer)
	executive.ExecuteBlock(block)
}
