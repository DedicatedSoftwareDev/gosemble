package aura

import (
	"bytes"
	"reflect"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants"
	"github.com/LimeChain/gosemble/constants/aura"
	"github.com/LimeChain/gosemble/constants/timestamp"
	"github.com/LimeChain/gosemble/frame/system"
	"github.com/LimeChain/gosemble/primitives/hashing"
	"github.com/LimeChain/gosemble/primitives/log"
	"github.com/LimeChain/gosemble/primitives/storage"
	"github.com/LimeChain/gosemble/primitives/types"
	"github.com/LimeChain/gosemble/utils"
)

type Slot = sc.U64

// Authorities returns current set of AuRa (Authority Round) authorities.
// Returns a pointer-size of the SCALE-encoded set of authorities.
func Authorities() int64 {
	auraHash := hashing.Twox128(constants.KeyAura)
	authoritiesHash := hashing.Twox128(constants.KeyAuthorities)

	authorities := storage.Get(append(auraHash, authoritiesHash...))

	if !authorities.HasValue {
		return utils.BytesToOffsetAndSize([]byte{0})
	}

	return utils.BytesToOffsetAndSize(sc.SequenceU8ToBytes(authorities.Value))
}

// SlotDuration returns the slot duration for AuRa.
// Returns a pointer-size of the SCALE-encoded slot duration
func SlotDuration() int64 {
	slotDuration := sc.U64(slotDuration())
	return utils.BytesToOffsetAndSize(slotDuration.Bytes())
}

func OnGenesisSession() {
	// TODO: implement once Session module is added
}

func OnNewSession() {
	// TODO: implement once Session module is added
}

func OnTimestampSet(now sc.U64) {
	slotDuration := slotDuration()
	if slotDuration == 0 {
		log.Critical("Aura slot duration cannot be zero.")
	}

	timestampSlot := now / sc.U64(slotDuration)

	auraHash := hashing.Twox128(constants.KeyAura)
	currentSlotHash := hashing.Twox128(constants.KeyCurrentSlot)
	currentSlot := storage.GetDecode(append(auraHash, currentSlotHash...), sc.DecodeU64)

	if currentSlot != timestampSlot {
		log.Critical("Timestamp slot must match `CurrentSlot`")
	}
}

func currentSlotFromDigests() sc.Option[Slot] {
	digest := system.StorageGetDigest()

	for keyDigest, dig := range digest {
		if keyDigest == types.DigestTypePreRuntime {
			for _, digestItem := range dig {
				if reflect.DeepEqual(sc.FixedSequenceU8ToBytes(digestItem.Engine), aura.EngineId[:]) {
					buffer := &bytes.Buffer{}
					buffer.Write(sc.SequenceU8ToBytes(digestItem.Payload))

					return sc.NewOption[Slot](sc.DecodeU64(buffer))
				}
			}
		}
	}

	return sc.NewOption[Slot](nil)
}

func totalAuthorities() sc.Option[sc.U64] {
	auraHash := hashing.Twox128(constants.KeyAura)
	authoritiesHash := hashing.Twox128(constants.KeyAuthorities)

	// `Compact<u32>` is 5 bytes in maximum.
	data := [5]byte{}
	option := storage.Read(append(auraHash, authoritiesHash...), data[:], 0)

	if !option.HasValue {
		return sc.NewOption[sc.U64](nil)
	}

	length := option.Value
	if length > sc.U32(len(data)) {
		length = sc.U32(len(data))
	}

	buffer := &bytes.Buffer{}
	buffer.Write(data[:length])

	compact := sc.DecodeCompact(buffer)

	totalAuthorities := sc.U64(compact.ToBigInt().Uint64())

	return sc.NewOption[sc.U64](totalAuthorities)
}

func slotDuration() int {
	return timestamp.MinimumPeriod * 2
}
