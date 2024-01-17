package inherent

import (
	tsc "github.com/LimeChain/gosemble/constants/timestamp"
	"github.com/LimeChain/gosemble/execution/types"
	"github.com/LimeChain/gosemble/frame/timestamp"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

func CheckExtrinsics(data primitives.InherentData, block types.Block) primitives.CheckInherentsResult {
	result := primitives.NewCheckInherentsResult()

	for _, extrinsic := range block.Extrinsics {
		// Inherents are before any other extrinsics.
		// And signed extrinsics are not inherents.
		if extrinsic.IsSigned() {
			break
		}

		isInherent := false
		call := extrinsic.Function

		switch call.ModuleIndex() {
		case tsc.ModuleIndex:
			switch call.FunctionIndex() {
			case tsc.FunctionSetIndex:
				isInherent = true
				err := timestamp.CheckInherent(call.Args(), data)
				if err != nil {
					err := result.PutError(tsc.InherentIdentifier, err.(primitives.IsFatalError))
					if err != nil {
						panic(err)
					}

					if result.FatalError {
						return result
					}
				}
			}
		}

		// Inherents are before any other extrinsics.
		// No module marked it as inherent thus it is not.
		if !isInherent {
			break
		}
	}

	// TODO: go through all required pallets with required inherents

	return result
}
