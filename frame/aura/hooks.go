package aura

import (
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants"
	"github.com/LimeChain/gosemble/primitives/hashing"
	"github.com/LimeChain/gosemble/primitives/log"
	"github.com/LimeChain/gosemble/primitives/storage"
	"github.com/LimeChain/gosemble/primitives/types"
)

func OnInitialize() types.Weight {
	slot := currentSlotFromDigests()

	if slot.HasValue {
		newSlot := slot.Value

		auraHash := hashing.Twox128(constants.KeyAura)
		currentSlotHash := hashing.Twox128(constants.KeyCurrentSlot)

		currentSlot := storage.GetDecode(append(auraHash, currentSlotHash...), sc.DecodeU64)

		if currentSlot >= newSlot {
			log.Critical("Slot must increase")
		}

		storage.Set(append(auraHash, currentSlotHash...), newSlot.Bytes())

		totalAuthorities := totalAuthorities()
		if totalAuthorities.HasValue {
			_ = currentSlot % totalAuthorities.Value

			// TODO: implement once  Session module is added
			/*
				if T::DisabledValidators::is_disabled(authority_index as u32) {
							panic!(
								"Validator with index {:?} is disabled and should not be attempting to author blocks.",
								authority_index,
							);
						}
			*/
		}

		return constants.DbWeight.ReadsWrites(2, 1)
	} else {
		return constants.DbWeight.Reads(1)
	}
}
