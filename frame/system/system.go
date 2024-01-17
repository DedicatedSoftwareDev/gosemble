package system

import (
	"bytes"
	"fmt"
	"math"
	"reflect"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants"
	"github.com/LimeChain/gosemble/primitives/hashing"
	"github.com/LimeChain/gosemble/primitives/log"
	"github.com/LimeChain/gosemble/primitives/storage"
	"github.com/LimeChain/gosemble/primitives/trie"
	"github.com/LimeChain/gosemble/primitives/types"
)

func Finalize() types.Header {
	systemHash := hashing.Twox128(constants.KeySystem)

	StorageClearExecutionPhase()
	StorageClearAllExtrinsicsLength()

	blockNumber := StorageGetBlockNumber()
	parentHash := StorageGetParentHash()
	digest := StorageGetDigest()
	extrinsicCount := StorageGetExtrinsicCount(true)

	var extrinsics []byte
	extrinsicDataPrefixHash := append(systemHash, hashing.Twox128(constants.KeyExtrinsicData)...)

	for i := 0; i < int(extrinsicCount); i++ {
		sci := sc.U32(i)
		hashIndex := hashing.Twox64(sci.Bytes())

		extrinsicDataHashIndexHash := append(extrinsicDataPrefixHash, hashIndex...)
		extrinsic := storage.TakeBytes(append(extrinsicDataHashIndexHash, sci.Bytes()...))

		extrinsics = append(extrinsics, extrinsic...)
	}

	buf := &bytes.Buffer{}
	extrinsicsRootBytes := trie.Blake2256OrderedRoot(append(sc.ToCompact(uint64(extrinsicCount)).Bytes(), extrinsics...), constants.StorageVersion)
	buf.Write(extrinsicsRootBytes)
	extrinsicsRoot := types.DecodeH256(buf)
	buf.Reset()

	// saturating_sub
	toRemove := blockNumber - constants.BlockHashCount - 1
	if toRemove > blockNumber {
		toRemove = 0
	}

	if toRemove != 0 {
		blockNumHash := hashing.Twox64(toRemove.Bytes())
		blockNumKey := append(systemHash, hashing.Twox128(constants.KeyBlockHash)...)
		blockNumKey = append(blockNumKey, blockNumHash...)
		blockNumKey = append(blockNumKey, toRemove.Bytes()...)

		storage.Clear(blockNumKey)
	}

	storageRootBytes := storage.Root(int32(constants.RuntimeVersion.StateVersion))
	buf.Write(storageRootBytes)
	storageRoot := types.DecodeH256(buf)
	buf.Reset()

	return types.Header{
		ExtrinsicsRoot: extrinsicsRoot,
		StateRoot:      storageRoot,
		ParentHash:     parentHash,
		Number:         blockNumber,
		Digest:         digest,
	}
}

func Initialize(blockNumber types.BlockNumber, parentHash types.Blake2bHash, digest types.Digest) {
	StorageSetExecutionPhase(types.NewExtrinsicPhaseInitialization())
	StorageSetExtrinsicIndex(sc.U32(0))
	StorageSetBlockNumber(blockNumber)
	StorageSetDigest(digest)
	StorageSetParentHash(parentHash)
	StorageSetBlockHash(blockNumber-1, parentHash)
	StorageClearBlockWeight()
}

func NoteFinishedInitialize() {
	StorageSetExecutionPhase(types.NewExtrinsicPhaseApply(sc.U32(0)))
}

func NoteFinishedExtrinsics() {
	extrinsicIndex := StorageGetExtrinsicIndex(true)
	StorageSetExtrinsicCount(extrinsicIndex)
	StorageSetExecutionPhase(types.NewExtrinsicPhaseFinalization())
}

// Note what the extrinsic data of the current extrinsic index is.
//
// This is required to be called before applying an extrinsic. The data will used
// in [`finalize`] to calculate the correct extrinsics root.
func NoteExtrinsic(encodedExt []byte) {
	keySystemHash := hashing.Twox128(constants.KeySystem)
	keyExtrinsicData := hashing.Twox128(constants.KeyExtrinsicData)

	keyExtrinsicDataPrefixHash := append(keySystemHash, keyExtrinsicData...)
	extrinsicIndex := StorageGetExtrinsicIndex(false)

	hashIndex := hashing.Twox64(extrinsicIndex.Bytes())

	keySystemExtrinsicDataHashIndex := append(keyExtrinsicDataPrefixHash, hashIndex...)
	storage.Set(append(keySystemExtrinsicDataHashIndex, extrinsicIndex.Bytes()...), sc.BytesToSequenceU8(encodedExt).Bytes())
}

// NoteAppliedExtrinsic - To be called immediately after an extrinsic has been applied.
//
// Emits an `ExtrinsicSuccess` or `ExtrinsicFailed` event depending on the outcome.
// The emitted event contains the post-dispatch corrected weight including
// the base-weight for its dispatch class.
func NoteAppliedExtrinsic(r *types.DispatchResultWithPostInfo[types.PostDispatchInfo], info types.DispatchInfo) {
	baseWeight := DefaultBlockWeights().Get(info.Class).BaseExtrinsic
	info.Weight = types.ExtractActualWeight(r, &info).SaturatingAdd(baseWeight)
	info.PaysFee = types.ExtractActualPaysFee(r, &info)

	if r.HasError {
		log.Trace(fmt.Sprintf("Extrinsic failed at block(%d): {%v}", StorageGetBlockNumber(), r.Err))
		DepositEvent(NewEventExtrinsicFailed(r.Err.Error, info))
	} else {
		DepositEvent(NewEventExtrinsicSuccess(info))
	}

	nextExtrinsicIndex := StorageGetExtrinsicIndex(false) + sc.U32(1)

	keySystemHash := hashing.Twox128(constants.KeySystem)

	StorageSetExtrinsicIndex(nextExtrinsicIndex)

	keyExecutionPhaseHash := hashing.Twox128(constants.KeyExecutionPhase)
	storage.Set(append(keySystemHash, keyExecutionPhaseHash...), types.NewExtrinsicPhaseApply(nextExtrinsicIndex).Bytes())
}

func Mutate(who types.Address32, f func(who *types.AccountInfo) sc.Result[sc.Encodable]) sc.Result[sc.Encodable] {
	accountInfo := StorageGetAccount(who.FixedSequence)

	result := f(&accountInfo)
	if !result.HasError {
		systemHash := hashing.Twox128(constants.KeySystem)
		accountHash := hashing.Twox128(constants.KeyAccount)

		whoBytes := sc.FixedSequenceU8ToBytes(who.FixedSequence)

		key := append(systemHash, accountHash...)
		key = append(key, hashing.Blake128(whoBytes)...)
		key = append(key, whoBytes...)

		storage.Set(key, accountInfo.Bytes())
	}

	return result
}

func TryMutateExists(who types.Address32, f func(who *types.AccountData) sc.Result[sc.Encodable]) sc.Result[sc.Encodable] {
	account := StorageGetAccount(who.FixedSequence)
	wasProviding := false
	if !reflect.DeepEqual(account.Data, types.AccountData{}) {
		wasProviding = true
	}

	someData := &types.AccountData{}
	if wasProviding {
		someData = &account.Data
	}

	result := f(someData)
	if result.HasError {
		return result
	}

	isProviding := !reflect.DeepEqual(someData, types.AccountData{})

	if !wasProviding && isProviding {
		incProviders(who)
	} else if wasProviding && !isProviding {
		status, err := decProviders(who)
		if err != nil {
			return sc.Result[sc.Encodable]{
				HasError: true,
				Value:    err,
			}
		}
		if status == types.DecRefStatusExists {
			return result
		}
	} else if !wasProviding && !isProviding {
		return result
	}

	Mutate(who, func(a *types.AccountInfo) sc.Result[sc.Encodable] {
		if someData != nil {
			a.Data = *someData
		} else {
			a.Data = types.AccountData{}
		}

		return sc.Result[sc.Encodable]{}
	})

	return result
}

func AccountTryMutateExists(who types.Address32, f func(who *types.AccountInfo) sc.Result[sc.Encodable]) sc.Result[sc.Encodable] {
	account := StorageGetAccount(who.FixedSequence)

	result := f(&account)

	if !result.HasError {
		StorageSetAccount(who.FixedSequence, account)
	}

	return result
}

func incProviders(who types.Address32) types.IncRefStatus {
	result := Mutate(who, func(a *types.AccountInfo) sc.Result[sc.Encodable] {
		if a.Providers == 0 && a.Sufficients == 0 {
			a.Providers = 1
			onCreatedAccount(who)

			return sc.Result[sc.Encodable]{
				HasError: false,
				Value:    types.IncRefStatusCreated,
			}
		} else {
			// saturating_add
			newProviders := a.Providers + 1
			if newProviders < a.Providers {
				newProviders = math.MaxUint32
			}

			return sc.Result[sc.Encodable]{
				HasError: false,
				Value:    types.IncRefStatusExisted,
			}
		}
	})

	return result.Value.(types.IncRefStatus)
}

func decProviders(who types.Address32) (types.DecRefStatus, types.DispatchError) {
	result := AccountTryMutateExists(who, func(account *types.AccountInfo) sc.Result[sc.Encodable] {
		if account.Providers == 0 {
			log.Warn("Logic error: Unexpected underflow in reducing provider")

			account.Providers = 1
		}

		if account.Providers == 1 && account.Consumers == 0 && account.Sufficients == 0 {
			return sc.Result[sc.Encodable]{
				HasError: false,
				Value:    types.DecRefStatusReaped,
			}
		}

		if account.Providers == 1 && account.Consumers > 0 {
			return sc.Result[sc.Encodable]{
				HasError: true,
				Value:    types.NewDispatchErrorConsumerRemaining(),
			}
		}

		account.Providers -= 1
		return sc.Result[sc.Encodable]{
			HasError: false,
			Value:    types.DecRefStatusExists,
		}
	})

	if result.HasError {
		return sc.U8(0), result.Value.(types.DispatchError)
	}

	return result.Value.(types.DecRefStatus), nil
}

func CanDecProviders(who types.Address32) bool {
	acc := StorageGetAccount(who.FixedSequence)

	return acc.Consumers == 0 || acc.Providers > 1
}

// RegisterExtraWeightUnchecked - Inform the system pallet of some additional weight that should be accounted for, in the
// current block.
//
// NOTE: use with extra care; this function is made public only be used for certain pallets
// that need it. A runtime that does not have dynamic calls should never need this and should
// stick to static weights. A typical use case for this is inner calls or smart contract calls.
// Furthermore, it only makes sense to use this when it is presumably  _cheap_ to provide the
// argument `weight`; In other words, if this function is to be used to account for some
// unknown, user provided call's weight, it would only make sense to use it if you are sure you
// can rapidly compute the weight of the inner call.
//
// Even more dangerous is to note that this function does NOT take any action, if the new sum
// of block weight is more than the block weight limit. This is what the _unchecked_.
//
// Another potential use-case could be for the `on_initialize` and `on_finalize` hooks.
func RegisterExtraWeightUnchecked(weight types.Weight, class types.DispatchClass) {
	currentWeight := StorageGetBlockWeight()
	currentWeight.Accrue(weight, class)
	StorageSetBlockWeight(currentWeight)
}

func ResetEvents() {
	StorageClearEvents()
	StorageClearEventCount()
	limit := sc.U32(math.MaxUint32)
	StorageClearEventTopics(limit)
}
