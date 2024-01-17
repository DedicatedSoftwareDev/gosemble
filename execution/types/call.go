package types

import (
	"bytes"
	"fmt"

	sc "github.com/LimeChain/goscale"

	"github.com/LimeChain/gosemble/config"
	"github.com/LimeChain/gosemble/primitives/log"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

func DecodeCall(buffer *bytes.Buffer) primitives.Call {
	moduleIndex := sc.DecodeU8(buffer)
	functionIndex := sc.DecodeU8(buffer)

	module, ok := config.Modules[moduleIndex]
	if !ok {
		log.Critical(fmt.Sprintf("module with index [%d] not found", moduleIndex))
	}

	function, ok := module.Functions()[functionIndex]
	if !ok {
		log.Critical(fmt.Sprintf("function index [%d] for module [%d] not found", functionIndex, moduleIndex))
	}

	function = function.DecodeArgs(buffer)

	return function
}
