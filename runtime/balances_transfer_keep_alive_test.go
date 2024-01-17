package main

import (
	"bytes"
	"math/big"
	"testing"

	gossamertypes "github.com/ChainSafe/gossamer/dot/types"
	"github.com/ChainSafe/gossamer/lib/common"
	"github.com/ChainSafe/gossamer/pkg/scale"
	"github.com/LimeChain/gosemble/constants"
	primitives "github.com/LimeChain/gosemble/primitives/types"
	cscale "github.com/centrifuge/go-substrate-rpc-client/v4/scale"
	"github.com/centrifuge/go-substrate-rpc-client/v4/signature"
	ctypes "github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/stretchr/testify/assert"
)

func Test_Balances_TransferKeepAlive_Success(t *testing.T) {
	rt, storage := newTestRuntime(t)
	runtimeVersion, err := rt.Version()
	assert.NoError(t, err)

	metadata := runtimeMetadata(t, rt)

	bob, err := ctypes.NewMultiAddressFromHexAccountID(
		"0x90b5ab205c6974c9ea841be688864633dc9ca8a357843eeacf2314649965fe22")
	assert.NoError(t, err)

	transferAmount := big.NewInt(0).SetUint64(constants.Dollar)

	call, err := ctypes.NewCall(metadata, "Balances.transfer_keep_alive", bob, ctypes.NewUCompact(transferAmount))
	assert.NoError(t, err)

	// Create the extrinsic
	ext := ctypes.NewExtrinsic(call)
	o := ctypes.SignatureOptions{
		BlockHash:          ctypes.Hash(parentHash),
		Era:                ctypes.ExtrinsicEra{IsImmortalEra: true},
		GenesisHash:        ctypes.Hash(parentHash),
		Nonce:              ctypes.NewUCompactFromUInt(0),
		SpecVersion:        ctypes.U32(runtimeVersion.SpecVersion),
		Tip:                ctypes.NewUCompactFromUInt(0),
		TransactionVersion: ctypes.U32(runtimeVersion.TransactionVersion),
	}

	// Set Account Info
	balance, e := big.NewInt(0).SetString("500000000000000", 10)
	assert.True(t, e)

	keyStorageAccountAlice, aliceAccountInfo := setStorageAccountInfo(t, storage, signature.TestKeyringPairAlice.PublicKey, balance, 0)

	// Sign the transaction using Alice's default account
	err = ext.Sign(signature.TestKeyringPairAlice, o)
	assert.NoError(t, err)

	extEnc := bytes.Buffer{}
	encoder := cscale.NewEncoder(&extEnc)
	err = ext.Encode(*encoder)
	assert.NoError(t, err)

	header := gossamertypes.NewHeader(parentHash, stateRoot, extrinsicsRoot, blockNumber, gossamertypes.NewDigest())
	encodedHeader, err := scale.Marshal(*header)
	assert.NoError(t, err)

	_, err = rt.Exec("Core_initialize_block", encodedHeader)
	assert.NoError(t, err)

	queryInfo := getQueryInfo(t, rt, extEnc.Bytes())

	res, err := rt.Exec("BlockBuilder_apply_extrinsic", extEnc.Bytes())
	assert.NoError(t, err)

	assert.Equal(t,
		primitives.NewApplyExtrinsicResult(primitives.NewDispatchOutcome(nil)).Bytes(),
		res,
	)

	bobHash, _ := common.Blake2b128(bob.AsID[:])
	keyStorageAccountBob := append(keySystemHash, keyAccountHash...)
	keyStorageAccountBob = append(keyStorageAccountBob, bobHash...)
	keyStorageAccountBob = append(keyStorageAccountBob, bob.AsID[:]...)
	bytesStorageBob := (*storage).Get(keyStorageAccountBob)

	expectedBobAccountInfo := gossamertypes.AccountInfo{
		Nonce:       0,
		Consumers:   0,
		Producers:   1,
		Sufficients: 0,
		Data: gossamertypes.AccountData{
			Free:       scale.MustNewUint128(transferAmount),
			Reserved:   scale.MustNewUint128(big.NewInt(0)),
			MiscFrozen: scale.MustNewUint128(big.NewInt(0)),
			FreeFrozen: scale.MustNewUint128(big.NewInt(0)),
		},
	}

	bobAccountInfo := gossamertypes.AccountInfo{}

	err = scale.Unmarshal(bytesStorageBob, &bobAccountInfo)
	assert.NoError(t, err)

	assert.Equal(t, expectedBobAccountInfo, bobAccountInfo)

	expectedAliceFreeBalance := big.NewInt(0).Sub(
		balance,
		big.NewInt(0).
			Add(transferAmount, queryInfo.PartialFee.ToBigInt()))
	expectedAliceAccountInfo := gossamertypes.AccountInfo{
		Nonce:       1,
		Consumers:   0,
		Producers:   0,
		Sufficients: 0,
		Data: gossamertypes.AccountData{
			Free:       scale.MustNewUint128(expectedAliceFreeBalance),
			Reserved:   scale.MustNewUint128(big.NewInt(0)),
			MiscFrozen: scale.MustNewUint128(big.NewInt(0)),
			FreeFrozen: scale.MustNewUint128(big.NewInt(0)),
		},
	}

	bytesAliceStorage := (*storage).Get(keyStorageAccountAlice)
	err = scale.Unmarshal(bytesAliceStorage, &aliceAccountInfo)
	assert.NoError(t, err)

	assert.Equal(t, expectedAliceAccountInfo, aliceAccountInfo)
}
