package main

import (
	"bytes"
	"math/big"
	"testing"

	gossamertypes "github.com/ChainSafe/gossamer/dot/types"
	"github.com/ChainSafe/gossamer/lib/common"
	"github.com/ChainSafe/gossamer/pkg/scale"
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants/balances"
	"github.com/LimeChain/gosemble/frame/balances/errors"
	primitives "github.com/LimeChain/gosemble/primitives/types"
	cscale "github.com/centrifuge/go-substrate-rpc-client/v4/scale"
	"github.com/centrifuge/go-substrate-rpc-client/v4/signature"
	ctypes "github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/stretchr/testify/assert"
)

func Test_Balances_TransferAll_Success_AllowDeath(t *testing.T) {
	rt, storage := newTestRuntime(t)
	runtimeVersion, err := rt.Version()
	assert.NoError(t, err)

	metadata := runtimeMetadata(t, rt)

	bob, err := ctypes.NewMultiAddressFromHexAccountID(
		"0x90b5ab205c6974c9ea841be688864633dc9ca8a357843eeacf2314649965fe22")
	assert.NoError(t, err)

	call, err := ctypes.NewCall(metadata, "Balances.transfer_all", bob, ctypes.NewBool(false))
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
			Free:       scale.MustNewUint128(big.NewInt(0).Sub(balance, queryInfo.PartialFee.ToBigInt())),
			Reserved:   scale.MustNewUint128(big.NewInt(0)),
			MiscFrozen: scale.MustNewUint128(big.NewInt(0)),
			FreeFrozen: scale.MustNewUint128(big.NewInt(0)),
		},
	}

	bobAccountInfo := gossamertypes.AccountInfo{}

	err = scale.Unmarshal(bytesStorageBob, &bobAccountInfo)
	assert.NoError(t, err)

	assert.Equal(t, expectedBobAccountInfo, bobAccountInfo)

	expectedAliceAccountInfo := gossamertypes.AccountInfo{
		Nonce:       1,
		Consumers:   0,
		Producers:   0,
		Sufficients: 0,
		Data: gossamertypes.AccountData{
			Free:       scale.MustNewUint128(big.NewInt(0)),
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

func Test_Balances_TransferAll_Success_KeepAlive(t *testing.T) {
	rt, storage := newTestRuntime(t)
	runtimeVersion, err := rt.Version()
	assert.NoError(t, err)

	metadata := runtimeMetadata(t, rt)

	bob, err := ctypes.NewMultiAddressFromHexAccountID(
		"0x90b5ab205c6974c9ea841be688864633dc9ca8a357843eeacf2314649965fe22")
	assert.NoError(t, err)

	call, err := ctypes.NewCall(metadata, "Balances.transfer_all", bob, ctypes.NewBool(true))
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

	setStorageAccountInfo(t, storage, signature.TestKeyringPairAlice.PublicKey, balance, 0)

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

	res, err := rt.Exec("BlockBuilder_apply_extrinsic", extEnc.Bytes())
	assert.NoError(t, err)

	// TODO: remove once tx payments are implemented
	expectedResult :=
		primitives.NewApplyExtrinsicResult(
			primitives.NewDispatchOutcome(
				primitives.NewDispatchErrorModule(
					primitives.CustomModuleError{
						Index: balances.ModuleIndex,
						Error: sc.U32(errors.ErrorKeepAlive),
					})))

	assert.Equal(t,
		expectedResult.Bytes(),
		res,
	)

	// TODO: Uncomment once tx payments are implemented, this will be successfully executed,
	// for now it fails due to nothing reserved in account executor
	//assert.Equal(t,
	//	primitives.NewApplyExtrinsicResult(primitives.NewDispatchOutcome(nil)).Bytes(),
	//	res,
	//)

	//bobHash, _ := common.Blake2b128(bob.AsID[:])
	//keyStorageAccountBob := append(keySystemHash, keyAccountHash...)
	//keyStorageAccountBob = append(keyStorageAccountBob, bobHash...)
	//keyStorageAccountBob = append(keyStorageAccountBob, bob.AsID[:]...)
	//bytesStorageBob := storage.Get(keyStorageAccountBob)
	//
	//expectedBobAccountInfo := gossamertypes.AccountInfo{
	//	Nonce:       0,
	//	Consumers:   0,
	//	Producers:   1,
	//	Sufficients: 0,
	//	Data: gossamertypes.AccountData{
	//		Free:       scale.MustNewUint128(mockBalance),
	//		Reserved:   scale.MustNewUint128(big.NewInt(0)),
	//		MiscFrozen: scale.MustNewUint128(big.NewInt(0)),
	//		FreeFrozen: scale.MustNewUint128(big.NewInt(0)),
	//	},
	//}
	//
	//bobAccountInfo := gossamertypes.AccountInfo{}
	//
	//err = scale.Unmarshal(bytesStorageBob, &bobAccountInfo)
	//assert.NoError(t, err)
	//
	//assert.Equal(t, expectedBobAccountInfo, bobAccountInfo)
	//
	//expectedAliceAccountInfo := gossamertypes.AccountInfo{
	//	Nonce:       1,
	//	Consumers:   0,
	//	Producers:   0,
	//	Sufficients: 0,
	//	Data: gossamertypes.AccountData{
	//		Free:       scale.MustNewUint128(big.NewInt(0)),
	//		Reserved:   scale.MustNewUint128(big.NewInt(0)),
	//		MiscFrozen: scale.MustNewUint128(big.NewInt(0)),
	//		FreeFrozen: scale.MustNewUint128(big.NewInt(0)),
	//	},
	//}
	//
	//bytesAliceStorage := storage.Get(keyStorageAccountAlice)
	//err = scale.Unmarshal(bytesAliceStorage, &aliceAccountInfo)
	//assert.NoError(t, err)
	//
	//assert.Equal(t, expectedAliceAccountInfo, aliceAccountInfo)
}
