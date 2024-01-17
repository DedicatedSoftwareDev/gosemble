package main

import (
	"bytes"
	"testing"
	"time"

	gossamertypes "github.com/ChainSafe/gossamer/dot/types"
	"github.com/ChainSafe/gossamer/lib/common"
	"github.com/ChainSafe/gossamer/pkg/scale"
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants"
	"github.com/LimeChain/gosemble/constants/aura"
	"github.com/LimeChain/gosemble/execution/types"
	timestamp "github.com/LimeChain/gosemble/frame/timestamp/dispatchables"
	primitives "github.com/LimeChain/gosemble/primitives/types"
	"github.com/stretchr/testify/assert"
)

func Test_BlockExecution(t *testing.T) {
	// core.InitializeBlock
	// blockBuilder.InherentExtrinsics
	// blockBuilder.ApplyExtrinsics
	// blockBuilder.FinalizeBlock

	storageRoot := common.MustHexToHash("0x152062ceba215bcc9c8fa2acc30093d4ba1eb5a19a2ef57aed5b1ecb4ed52812") // Depends on timestamp
	time := time.Date(2023, time.January, 2, 3, 4, 5, 6, time.UTC)

	expectedStorageDigest := gossamertypes.NewDigest()
	digest := gossamertypes.NewDigest()

	rt, storage := newTestRuntime(t)

	bytesSlotDuration, err := rt.Exec("AuraApi_slot_duration", []byte{})
	assert.NoError(t, err)

	buffer := &bytes.Buffer{}
	buffer.Write(bytesSlotDuration)

	slotDuration := sc.DecodeU64(buffer)
	buffer.Reset()

	slot := sc.U64(time.UnixMilli()) / slotDuration

	preRuntimeDigest := gossamertypes.PreRuntimeDigest{
		ConsensusEngineID: aura.EngineId,
		Data:              slot.Bytes(),
	}
	assert.NoError(t, digest.Add(preRuntimeDigest))
	assert.NoError(t, expectedStorageDigest.Add(preRuntimeDigest))

	header := gossamertypes.NewHeader(parentHash, storageRoot, extrinsicsRoot, blockNumber, digest)

	encodedHeader, err := scale.Marshal(*header)
	assert.NoError(t, err)

	_, err = rt.Exec("Core_initialize_block", encodedHeader)
	assert.NoError(t, err)

	lrui := primitives.LastRuntimeUpgradeInfo{
		SpecVersion: sc.ToCompact(constants.SpecVersion),
		SpecName:    constants.SpecName,
	}
	assert.Equal(t, lrui.Bytes(), (*storage).Get(append(keySystemHash, keyLastRuntime...)))

	encExtrinsicIndex0, _ := scale.Marshal(uint32(0))
	assert.Equal(t, encExtrinsicIndex0, (*storage).Get(constants.KeyExtrinsicIndex))

	expectedExecutionPhase := primitives.NewExtrinsicPhaseApply(sc.U32(0))
	assert.Equal(t, expectedExecutionPhase.Bytes(), (*storage).Get(append(keySystemHash, keyExecutionPhaseHash...)))

	encBlockNumber, _ := scale.Marshal(uint32(blockNumber))
	assert.Equal(t, encBlockNumber, (*storage).Get(append(keySystemHash, keyNumberHash...)))

	encExpectedDigest, err := scale.Marshal(expectedStorageDigest)
	assert.NoError(t, err)

	assert.Equal(t, encExpectedDigest, (*storage).Get(append(keySystemHash, keyDigestHash...)))
	assert.Equal(t, parentHash.ToBytes(), (*storage).Get(append(keySystemHash, keyParentHash...)))

	blockHashKey := append(keySystemHash, keyBlockHash...)
	encPrevBlock, _ := scale.Marshal(uint32(blockNumber - 1))
	numHash, err := common.Twox64(encPrevBlock)
	assert.NoError(t, err)

	blockHashKey = append(blockHashKey, numHash...)
	blockHashKey = append(blockHashKey, encPrevBlock...)
	assert.Equal(t, parentHash.ToBytes(), (*storage).Get(blockHashKey))

	idata := gossamertypes.NewInherentData()
	err = idata.SetInherent(gossamertypes.Timstap0, uint64(time.UnixMilli()))

	assert.NoError(t, err)

	call := timestamp.NewSetCall(sc.NewVaryingData(sc.ToCompact(time.UnixMilli())))

	expectedExtrinsic := types.NewUnsignedUncheckedExtrinsic(call)

	ienc, err := idata.Encode()
	assert.NoError(t, err)

	inherentExt, err := rt.Exec("BlockBuilder_inherent_extrinsics", ienc)
	assert.NoError(t, err)
	assert.NotNil(t, inherentExt)

	buffer.Write([]byte{inherentExt[0]})

	totalInherents := sc.DecodeCompact(buffer)
	assert.Equal(t, int64(1), totalInherents.ToBigInt().Int64())
	buffer.Reset()

	buffer.Write(inherentExt[1:])
	extrinsic := types.DecodeUncheckedExtrinsic(buffer)
	buffer.Reset()

	assert.Equal(t, expectedExtrinsic, extrinsic)

	applyResult, err := rt.Exec("BlockBuilder_apply_extrinsic", inherentExt[1:])
	assert.NoError(t, err)

	assert.Equal(t,
		primitives.NewApplyExtrinsicResult(primitives.NewDispatchOutcome(nil)).Bytes(),
		applyResult,
	)

	bytesResult, err := rt.Exec("BlockBuilder_finalize_block", []byte{})
	assert.NoError(t, err)

	resultHeader := gossamertypes.NewEmptyHeader()
	assert.NoError(t, scale.Unmarshal(bytesResult, resultHeader))
	resultHeader.Hash() // Call this to be set, otherwise structs do not match...

	assert.Equal(t, header, resultHeader)

	assert.Equal(t, []byte(nil), (*storage).Get(append(keyTimestampHash, keyTimestampDidUpdate...)))
	assert.Equal(t, sc.U64(time.UnixMilli()).Bytes(), (*storage).Get(append(keyTimestampHash, keyTimestampNowHash...)))

	assert.Equal(t, []byte(nil), (*storage).Get(constants.KeyExtrinsicIndex))
	assert.Equal(t, []byte(nil), (*storage).Get(append(keySystemHash, keyExecutionPhaseHash...)))
	assert.Equal(t, []byte(nil), (*storage).Get(append(keySystemHash, keyAllExtrinsicsLenHash...)))
	assert.Equal(t, []byte(nil), (*storage).Get(append(keySystemHash, keyExtrinsicCountHash...)))

	assert.Equal(t, parentHash.ToBytes(), (*storage).Get(append(keySystemHash, keyParentHash...)))
	assert.Equal(t, encExpectedDigest, (*storage).Get(append(keySystemHash, keyDigestHash...)))
	assert.Equal(t, encBlockNumber, (*storage).Get(append(keySystemHash, keyNumberHash...)))

	assert.Equal(t, slot.Bytes(), (*storage).Get(append(keyAuraHash, keyCurrentSlotHash...)))
}

func Test_ExecuteBlock(t *testing.T) {
	// blockBuilder.Inherent_Extrinsics
	// blockBuilder.ExecuteBlock

	storageRoot := common.MustHexToHash("0x152062ceba215bcc9c8fa2acc30093d4ba1eb5a19a2ef57aed5b1ecb4ed52812") // Depends on timestamp
	time := time.Date(2023, time.January, 2, 3, 4, 5, 6, time.UTC)

	rt, storage := newTestRuntime(t)

	bytesSlotDuration, err := rt.Exec("AuraApi_slot_duration", []byte{})
	assert.NoError(t, err)

	buffer := &bytes.Buffer{}
	buffer.Write(bytesSlotDuration)

	slotDuration := sc.DecodeU64(buffer)
	buffer.Reset()

	slot := sc.U64(time.UnixMilli()) / slotDuration

	preRuntimeDigest := gossamertypes.PreRuntimeDigest{
		ConsensusEngineID: aura.EngineId,
		Data:              slot.Bytes(),
	}

	idata := gossamertypes.NewInherentData()
	err = idata.SetInherent(gossamertypes.Timstap0, uint64(time.UnixMilli()))

	assert.NoError(t, err)

	ienc, err := idata.Encode()
	assert.NoError(t, err)

	call := timestamp.NewSetCall(sc.NewVaryingData(sc.ToCompact(time.UnixMilli())))

	expectedExtrinsic := types.NewUnsignedUncheckedExtrinsic(call)

	inherentExt, err := rt.Exec("BlockBuilder_inherent_extrinsics", ienc)
	assert.NoError(t, err)
	assert.NotNil(t, inherentExt)

	buffer.Write([]byte{inherentExt[0]})

	totalInherents := sc.DecodeCompact(buffer)
	assert.Equal(t, int64(1), totalInherents.ToBigInt().Int64())
	buffer.Reset()

	buffer.Write(inherentExt[1:])
	extrinsic := types.DecodeUncheckedExtrinsic(buffer)

	assert.Equal(t, expectedExtrinsic, extrinsic)

	var exts [][]byte
	err = scale.Unmarshal(inherentExt, &exts)
	assert.Nil(t, err)

	digest := gossamertypes.NewDigest()

	assert.NoError(t, err)
	assert.NoError(t, digest.Add(preRuntimeDigest))

	expectedStorageDigest, err := scale.Marshal(digest)
	assert.NoError(t, err)
	encBlockNumber, _ := scale.Marshal(uint32(blockNumber))

	header := gossamertypes.NewHeader(parentHash, storageRoot, extrinsicsRoot, blockNumber, digest)

	block := gossamertypes.Block{
		Header: *header,
		Body:   gossamertypes.BytesArrayToExtrinsics(exts),
	}

	encodedBlock, err := scale.Marshal(block)
	assert.Nil(t, err)

	_, err = rt.Exec("Core_execute_block", encodedBlock)
	assert.NoError(t, err)

	assert.Equal(t, []byte(nil), (*storage).Get(append(keyTimestampHash, keyTimestampDidUpdate...)))
	assert.Equal(t, sc.U64(time.UnixMilli()).Bytes(), (*storage).Get(append(keyTimestampHash, keyTimestampNowHash...)))

	assert.Equal(t, []byte(nil), (*storage).Get(constants.KeyExtrinsicIndex))
	assert.Equal(t, []byte(nil), (*storage).Get(append(keySystemHash, keyExecutionPhaseHash...)))
	assert.Equal(t, []byte(nil), (*storage).Get(append(keySystemHash, keyAllExtrinsicsLenHash...)))
	assert.Equal(t, []byte(nil), (*storage).Get(append(keySystemHash, keyExtrinsicCountHash...)))

	assert.Equal(t, parentHash.ToBytes(), (*storage).Get(append(keySystemHash, keyParentHash...)))
	assert.Equal(t, expectedStorageDigest, (*storage).Get(append(keySystemHash, keyDigestHash...)))
	assert.Equal(t, encBlockNumber, (*storage).Get(append(keySystemHash, keyNumberHash...)))

	assert.Equal(t, slot.Bytes(), (*storage).Get(append(keyAuraHash, keyCurrentSlotHash...)))
}
