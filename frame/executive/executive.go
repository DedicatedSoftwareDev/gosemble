package executive

import (
	"fmt"
	"reflect"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants"
	"github.com/LimeChain/gosemble/execution/extrinsic"
	"github.com/LimeChain/gosemble/execution/inherent"
	"github.com/LimeChain/gosemble/execution/types"
	"github.com/LimeChain/gosemble/frame/aura"
	"github.com/LimeChain/gosemble/frame/system"
	"github.com/LimeChain/gosemble/primitives/crypto"
	"github.com/LimeChain/gosemble/primitives/hashing"
	"github.com/LimeChain/gosemble/primitives/log"
	"github.com/LimeChain/gosemble/primitives/storage"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

// InitializeBlock initialises a block with the given header,
// starting the execution of a particular block.
func InitializeBlock(header primitives.Header) {
	log.Trace("init_block")
	system.ResetEvents()

	weight := primitives.WeightZero()
	if runtimeUpgrade() {
		weight = weight.SaturatingAdd(executeOnRuntimeUpgrade())
	}

	system.Initialize(header.Number, header.ParentHash, extractPreRuntimeDigest(header.Digest))

	// TODO: accumulate the weight from all pallets that have on_initialize
	weight = weight.SaturatingAdd(aura.OnInitialize())
	weight = weight.SaturatingAdd(system.DefaultBlockWeights().BaseBlock)
	// use in case of dynamic weight calculation
	system.RegisterExtraWeightUnchecked(weight, primitives.NewDispatchClassMandatory())

	system.NoteFinishedInitialize()
}

func ExecuteBlock(block types.Block) {
	log.Trace(fmt.Sprintf("execute_block %v", block.Header.Number))

	InitializeBlock(block.Header)

	initialChecks(block)

	crypto.ExtCryptoStartBatchVerify()
	executeExtrinsicsWithBookKeeping(block)
	if crypto.ExtCryptoFinishBatchVerify() != 1 {
		log.Critical("Signature verification failed")
	}

	finalChecks(&block.Header)
}

// ApplyExtrinsic applies extrinsic outside the block execution function.
//
// This doesn't attempt to validate anything regarding the block, but it builds a list of uxt
// hashes.
func ApplyExtrinsic(uxt types.UncheckedExtrinsic) (primitives.DispatchOutcome, primitives.TransactionValidityError) {
	encoded := uxt.Bytes()
	encodedLen := sc.ToCompact(len(encoded))

	log.Trace("apply_extrinsic")

	// Verify that the signature is good.
	xt, err := extrinsic.Unchecked(uxt).Check(primitives.DefaultAccountIdLookup())
	if err != nil {
		return primitives.DispatchOutcome{}, err
	}

	// We don't need to make sure to `note_extrinsic` only after we know it's going to be
	// executed to prevent it from leaking in storage since at this point, it will either
	// execute or panic (and revert storage changes).
	system.NoteExtrinsic(encoded)

	// AUDIT: Under no circumstances may this function panic from here onwards.

	// Decode parameters and dispatch
	dispatchInfo := primitives.GetDispatchInfo(xt.Function)
	log.Trace("get_dispatch_info: weight ref time " + dispatchInfo.Weight.RefTime.String())

	unsignedValidator := extrinsic.UnsignedValidatorForChecked{}
	res, err := extrinsic.Checked(xt).Apply(unsignedValidator, &dispatchInfo, encodedLen)
	if err != nil {
		return primitives.DispatchOutcome{}, err
	}

	// Mandatory(inherents) are not allowed to fail.
	//
	// The entire block should be discarded if an inherent fails to apply. Otherwise
	// it may open an attack vector.
	if res.HasError && dispatchInfo.Class.Is(primitives.DispatchClassMandatory) {
		return primitives.DispatchOutcome{}, primitives.NewTransactionValidityError(primitives.NewInvalidTransactionBadMandatory())
	}

	system.NoteAppliedExtrinsic(&res, dispatchInfo)

	if res.HasError {
		return primitives.NewDispatchOutcome(res.Err.Error), nil
	}

	return primitives.NewDispatchOutcome(nil), nil
}

// ValidateTransaction checks a given signed transaction for validity. This doesn't execute any
// side-effects; it merely checks whether the transaction would panic if it were included or
// not.
//
// Changes made to storage should be discarded.
func ValidateTransaction(source primitives.TransactionSource, uxt types.UncheckedExtrinsic, blockHash primitives.Blake2bHash) (ok primitives.ValidTransaction, err primitives.TransactionValidityError) {
	currentBlockNumber := system.StorageGetBlockNumber()
	system.Initialize(currentBlockNumber+1, blockHash, primitives.Digest{})

	log.Trace("validate_transaction")

	log.Trace("using_encoded")
	encodedLen := sc.ToCompact(len(uxt.Bytes()))

	log.Trace("check")
	xt, err := extrinsic.Unchecked(uxt).Check(primitives.DefaultAccountIdLookup())
	if err != nil {
		return ok, err
	}

	log.Trace("dispatch_info")
	dispatchInfo := primitives.GetDispatchInfo(xt.Function)

	if dispatchInfo.Class.Is(primitives.DispatchClassMandatory) {
		return ok, primitives.NewTransactionValidityError(primitives.NewInvalidTransactionMandatoryValidation())
	}

	log.Trace("validate")
	unsignedValidator := extrinsic.UnsignedValidatorForChecked{}
	return extrinsic.Checked(xt).Validate(unsignedValidator, source, &dispatchInfo, encodedLen)
}

func executeExtrinsicsWithBookKeeping(block types.Block) {
	for _, ext := range block.Extrinsics {
		_, err := ApplyExtrinsic(ext)
		if err != nil {
			log.Critical(string(err[0].Bytes()))
		}
	}

	system.NoteFinishedExtrinsics()

	IdleAndFinalizeHook(block.Header.Number)
}

func initialChecks(block types.Block) {
	log.Trace("initial_checks")

	header := block.Header
	blockNumber := header.Number

	if blockNumber > 0 {
		storageParentHash := system.StorageGetBlockHash(blockNumber - 1)

		if !reflect.DeepEqual(storageParentHash, header.ParentHash) {
			log.Critical("parent hash should be valid")
		}
	}

	inherentsAreFirst := inherent.EnsureInherentsAreFirst(block)

	if inherentsAreFirst >= 0 {
		log.Critical(fmt.Sprintf("invalid inherent position for extrinsic at index [%d]", inherentsAreFirst))
	}
}

func runtimeUpgrade() sc.Bool {
	systemHash := hashing.Twox128(constants.KeySystem)
	lastRuntimeUpgradeHash := hashing.Twox128(constants.KeyLastRuntimeUpgrade)

	keyLru := append(systemHash, lastRuntimeUpgradeHash...)
	lrupi := storage.GetDecode(keyLru, primitives.DecodeLastRuntimeUpgradeInfo)

	if constants.RuntimeVersion.SpecVersion > sc.U32(lrupi.SpecVersion.ToBigInt().Int64()) ||
		lrupi.SpecName != constants.RuntimeVersion.SpecName {

		valueLru := append(
			sc.ToCompact(constants.RuntimeVersion.SpecVersion).Bytes(),
			constants.RuntimeVersion.SpecName.Bytes()...)
		storage.Set(keyLru, valueLru)

		return true
	}

	return false
}

func extractPreRuntimeDigest(digest primitives.Digest) primitives.Digest {
	result := primitives.Digest{}
	for k, v := range digest {
		if k == primitives.DigestTypePreRuntime {
			result[k] = v
		}
	}

	return result
}

func finalChecks(header *primitives.Header) {
	newHeader := system.Finalize()

	if len(header.Digest) != len(newHeader.Digest) {
		log.Critical("Number of digest must match the calculated")
	}

	for key, digest := range header.Digest {
		otherDigest := newHeader.Digest[key]
		if !reflect.DeepEqual(digest, otherDigest) {
			log.Critical("digest item must match that calculated")
		}
	}

	if !reflect.DeepEqual(header.StateRoot, newHeader.StateRoot) {
		log.Critical("Storage root must match that calculated")
	}

	if !reflect.DeepEqual(header.ExtrinsicsRoot, newHeader.ExtrinsicsRoot) {
		log.Critical("Transaction trie must be valid")
	}
}

// Execute all `OnRuntimeUpgrade` of this runtime, and return the aggregate weight.
func executeOnRuntimeUpgrade() primitives.Weight {
	// TODO: ex: balances
	// call on_runtime_upgrade hook for all modules that implement it
	return onRuntimeUpgrade()
}
