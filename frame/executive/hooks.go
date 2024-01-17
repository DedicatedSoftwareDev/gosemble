package executive

import (
	"fmt"

	"github.com/LimeChain/gosemble/frame/system"
	"github.com/LimeChain/gosemble/frame/timestamp"
	"github.com/LimeChain/gosemble/primitives/log"
	"github.com/LimeChain/gosemble/primitives/types"
)

func IdleAndFinalizeHook(blockNumber types.BlockNumber) {
	weight := system.StorageGetBlockWeight()

	maxWeight := system.DefaultBlockWeights().MaxBlock
	remainingWeight := maxWeight.SaturatingSub(weight.Total())

	if remainingWeight.AllGt(types.WeightZero()) {
		// TODO: call on_idle hook for each pallet
		usedWeight := onIdle(blockNumber, remainingWeight)
		system.RegisterExtraWeightUnchecked(usedWeight, types.NewDispatchClassMandatory())
	}

	// Each pallet (babe, grandpa) has its own on_finalize that has to be implemented once it is supported
	timestamp.OnFinalize()
}

func onRuntimeUpgrade() types.Weight {
	return types.WeightFromParts(200, 0)
}

func onIdle(n types.BlockNumber, remainingWeight types.Weight) types.Weight {
	log.Trace(fmt.Sprintf("on_idle %v, %v)", n, remainingWeight))
	return types.WeightFromParts(175, 0)
}
