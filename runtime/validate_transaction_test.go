package main

import (
	"bytes"
	"math/big"
	"testing"
	"time"

	gossamertypes "github.com/ChainSafe/gossamer/dot/types"
	"github.com/ChainSafe/gossamer/pkg/scale"
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants"
	primitives "github.com/LimeChain/gosemble/primitives/types"
	cscale "github.com/centrifuge/go-substrate-rpc-client/v4/scale"
	"github.com/centrifuge/go-substrate-rpc-client/v4/signature"
	ctypes "github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/stretchr/testify/assert"
)

func Test_ValidateTransaction_Success(t *testing.T) {
	rt, storage := newTestRuntime(t)
	runtimeVersion, err := rt.Version()
	assert.NoError(t, err)

	metadata := runtimeMetadata(t, rt)

	// Set Account Info balance otherwise tx payment check will fail.
	balance, e := big.NewInt(0).SetString("500000000000000", 10)
	assert.True(t, e)

	setStorageAccountInfo(t, storage, signature.TestKeyringPairAlice.PublicKey, balance, 0)

	digest := gossamertypes.NewDigest()

	header := gossamertypes.NewHeader(parentHash, stateRoot, extrinsicsRoot, blockNumber, digest)
	encodedHeader, err := scale.Marshal(*header)
	assert.NoError(t, err)

	_, err = rt.Exec("Core_initialize_block", encodedHeader)
	assert.NoError(t, err)

	call, err := ctypes.NewCall(metadata, "System.remark", []byte{})
	assert.NoError(t, err)

	extrinsic := ctypes.NewExtrinsic(call)

	o := ctypes.SignatureOptions{
		BlockHash:          ctypes.Hash(parentHash),
		Era:                ctypes.ExtrinsicEra{IsImmortalEra: true},
		GenesisHash:        ctypes.Hash(parentHash),
		Nonce:              ctypes.NewUCompactFromUInt(0),
		SpecVersion:        ctypes.U32(runtimeVersion.SpecVersion),
		Tip:                ctypes.NewUCompactFromUInt(0),
		TransactionVersion: ctypes.U32(runtimeVersion.TransactionVersion),
	}

	// Sign the transaction using Alice's default account
	err = extrinsic.Sign(signature.TestKeyringPairAlice, o)
	assert.NoError(t, err)

	txSource := primitives.NewTransactionSourceExternal()
	blockHash := sc.BytesToFixedSequenceU8(parentHash.ToBytes())

	buffer := &bytes.Buffer{}
	txSource.Encode(buffer)

	encoder := cscale.NewEncoder(buffer)
	err = extrinsic.Encode(*encoder)
	assert.NoError(t, err)

	blockHash.Encode(buffer)

	encTransactionValidityResult, err := rt.Exec("TaggedTransactionQueue_validate_transaction", buffer.Bytes())
	assert.NoError(t, err)

	buffer.Reset()
	buffer.Write(encTransactionValidityResult)
	transactionValidityResult := primitives.DecodeTransactionValidityResult(buffer)

	assert.Equal(t, sc.Bool(true), transactionValidityResult.IsValidTransaction())
}

func Test_ValidateTransaction_InvalidModuleFunctionIndex(t *testing.T) {
	rt, _ := newTestRuntime(t)
	runtimeVersion, err := rt.Version()
	assert.NoError(t, err)

	metadata := runtimeMetadata(t, rt)

	digest := gossamertypes.NewDigest()

	header := gossamertypes.NewHeader(parentHash, stateRoot, extrinsicsRoot, blockNumber, digest)
	encodedHeader, err := scale.Marshal(*header)
	assert.NoError(t, err)

	_, err = rt.Exec("Core_initialize_block", encodedHeader)
	assert.NoError(t, err)

	call, err := ctypes.NewCall(metadata, "System.remark", []byte{})
	assert.NoError(t, err)

	// Change function/section index
	call.CallIndex.SectionIndex = 65

	extrinsic := ctypes.NewExtrinsic(call)

	o := ctypes.SignatureOptions{
		BlockHash:          ctypes.Hash(parentHash),
		Era:                ctypes.ExtrinsicEra{IsImmortalEra: true},
		GenesisHash:        ctypes.Hash(parentHash),
		Nonce:              ctypes.NewUCompactFromUInt(0),
		SpecVersion:        ctypes.U32(runtimeVersion.SpecVersion),
		Tip:                ctypes.NewUCompactFromUInt(0),
		TransactionVersion: ctypes.U32(runtimeVersion.TransactionVersion),
	}

	// Sign the transaction using Alice's default account
	err = extrinsic.Sign(signature.TestKeyringPairAlice, o)
	assert.NoError(t, err)

	txSource := primitives.NewTransactionSourceExternal()
	blockHash := sc.BytesToFixedSequenceU8(parentHash.ToBytes())

	buffer := &bytes.Buffer{}
	txSource.Encode(buffer)

	encoder := cscale.NewEncoder(buffer)
	err = extrinsic.Encode(*encoder)
	assert.NoError(t, err)

	blockHash.Encode(buffer)

	_, err = rt.Exec("TaggedTransactionQueue_validate_transaction", buffer.Bytes())
	assert.Error(t, err)
}

func Test_ValidateTransaction_StaleError_InvalidNonce(t *testing.T) {
	rt, storage := newTestRuntime(t)
	runtimeVersion, err := rt.Version()
	assert.NoError(t, err)

	metadata := runtimeMetadata(t, rt)

	setStorageAccountInfo(t, storage, signature.TestKeyringPairAlice.PublicKey, big.NewInt(5), 3)

	digest := gossamertypes.NewDigest()

	header := gossamertypes.NewHeader(parentHash, stateRoot, extrinsicsRoot, blockNumber, digest)
	encodedHeader, err := scale.Marshal(*header)
	assert.NoError(t, err)

	_, err = rt.Exec("Core_initialize_block", encodedHeader)
	assert.NoError(t, err)

	call, err := ctypes.NewCall(metadata, "System.remark", []byte{})
	assert.NoError(t, err)

	extrinsic := ctypes.NewExtrinsic(call)
	o := ctypes.SignatureOptions{
		BlockHash:          ctypes.Hash(parentHash),
		Era:                ctypes.ExtrinsicEra{IsImmortalEra: true},
		GenesisHash:        ctypes.Hash(parentHash),
		Nonce:              ctypes.NewUCompactFromUInt(2), // Invalid nonce
		SpecVersion:        ctypes.U32(runtimeVersion.SpecVersion),
		Tip:                ctypes.NewUCompactFromUInt(0),
		TransactionVersion: ctypes.U32(runtimeVersion.TransactionVersion),
	}

	// Sign the transaction using Alice's default account
	err = extrinsic.Sign(signature.TestKeyringPairAlice, o)
	assert.NoError(t, err)

	txSource := primitives.NewTransactionSourceExternal()
	blockHash := sc.BytesToFixedSequenceU8(parentHash.ToBytes())

	buffer := &bytes.Buffer{}
	txSource.Encode(buffer)

	encoder := cscale.NewEncoder(buffer)
	err = extrinsic.Encode(*encoder)
	assert.NoError(t, err)

	blockHash.Encode(buffer)

	encTransactionValidityResult, err := rt.Exec("TaggedTransactionQueue_validate_transaction", buffer.Bytes())
	assert.NoError(t, err)

	buffer.Reset()
	buffer.Write(encTransactionValidityResult)
	transactionValidityResult := primitives.DecodeTransactionValidityResult(buffer)

	assert.Equal(t,
		primitives.NewTransactionValidityResult(
			primitives.NewTransactionValidityError(
				primitives.NewInvalidTransactionStale(),
			),
		),
		transactionValidityResult,
	)
}

func Test_ValidateTransaction_ExhaustsResourcesError(t *testing.T) {
	rt, storage := newTestRuntime(t)
	runtimeVersion, err := rt.Version()
	assert.NoError(t, err)

	metadata := runtimeMetadata(t, rt)

	setStorageAccountInfo(t, storage, signature.TestKeyringPairAlice.PublicKey, big.NewInt(5), 0)

	digest := gossamertypes.NewDigest()

	header := gossamertypes.NewHeader(parentHash, stateRoot, extrinsicsRoot, blockNumber, digest)
	encodedHeader, err := scale.Marshal(*header)
	assert.NoError(t, err)

	_, err = rt.Exec("Core_initialize_block", encodedHeader)
	assert.NoError(t, err)

	// Append long args
	args := make([]byte, constants.FiveMbPerBlockPerExtrinsic)

	call, err := ctypes.NewCall(metadata, "System.remark", args)
	assert.NoError(t, err)

	extrinsic := ctypes.NewExtrinsic(call)
	o := ctypes.SignatureOptions{
		BlockHash:          ctypes.Hash(parentHash),
		Era:                ctypes.ExtrinsicEra{IsImmortalEra: true},
		GenesisHash:        ctypes.Hash(parentHash),
		Nonce:              ctypes.NewUCompactFromUInt(0),
		SpecVersion:        ctypes.U32(runtimeVersion.SpecVersion),
		Tip:                ctypes.NewUCompactFromUInt(0),
		TransactionVersion: ctypes.U32(runtimeVersion.TransactionVersion),
	}

	// Sign the transaction using Alice's default account
	err = extrinsic.Sign(signature.TestKeyringPairAlice, o)
	assert.NoError(t, err)

	txSource := primitives.NewTransactionSourceExternal()
	blockHash := sc.BytesToFixedSequenceU8(parentHash.ToBytes())

	buffer := &bytes.Buffer{}
	txSource.Encode(buffer)

	encoder := cscale.NewEncoder(buffer)
	err = extrinsic.Encode(*encoder)
	assert.NoError(t, err)

	blockHash.Encode(buffer)

	encTransactionValidityResult, err := rt.Exec("TaggedTransactionQueue_validate_transaction", buffer.Bytes())
	assert.NoError(t, err)

	buffer.Reset()
	buffer.Write(encTransactionValidityResult)
	transactionValidityResult := primitives.DecodeTransactionValidityResult(buffer)

	assert.Equal(t,
		primitives.NewTransactionValidityResult(
			primitives.NewTransactionValidityError(
				primitives.NewInvalidTransactionExhaustsResources(),
			),
		),
		transactionValidityResult,
	)
}

func Test_ValidateTransaction_Era(t *testing.T) {
	rt, storage := newTestRuntime(t)
	runtimeVersion, err := rt.Version()
	assert.NoError(t, err)

	metadata := runtimeMetadata(t, rt)

	// Set Account info due to check tx payment
	balance, e := big.NewInt(0).SetString("500000000000000", 10)
	assert.True(t, e)

	setStorageAccountInfo(t, storage, signature.TestKeyringPairAlice.PublicKey, balance, 0)

	digest := gossamertypes.NewDigest()

	header := gossamertypes.NewHeader(parentHash, stateRoot, extrinsicsRoot, blockNumber, digest)
	encodedHeader, err := scale.Marshal(*header)
	assert.NoError(t, err)

	_, err = rt.Exec("Core_initialize_block", encodedHeader)
	assert.NoError(t, err)

	setBlockNumber(t, storage, 16)

	call, err := ctypes.NewCall(metadata, "System.remark", []byte{})
	assert.NoError(t, err)

	extrinsic := ctypes.NewExtrinsic(call)

	o := ctypes.SignatureOptions{
		BlockHash: ctypes.Hash(parentHash),
		Era: ctypes.ExtrinsicEra{
			IsMortalEra: true,
			AsMortalEra: ctypes.MortalEra{
				First:  3, // Matched with period 16, current 256
				Second: 0,
			},
		},
		GenesisHash:        ctypes.Hash(parentHash),
		Nonce:              ctypes.NewUCompactFromUInt(0),
		SpecVersion:        ctypes.U32(runtimeVersion.SpecVersion),
		Tip:                ctypes.NewUCompactFromUInt(0),
		TransactionVersion: ctypes.U32(runtimeVersion.TransactionVersion),
	}

	// Sign the transaction using Alice's default account
	err = extrinsic.Sign(signature.TestKeyringPairAlice, o)
	assert.NoError(t, err)

	txSource := primitives.NewTransactionSourceExternal()
	blockHash := sc.BytesToFixedSequenceU8(parentHash.ToBytes())

	buffer := &bytes.Buffer{}
	txSource.Encode(buffer)

	encoder := cscale.NewEncoder(buffer)
	err = extrinsic.Encode(*encoder)
	assert.NoError(t, err)

	blockHash.Encode(buffer)

	encTransactionValidityResult, err := rt.Exec("TaggedTransactionQueue_validate_transaction", buffer.Bytes())
	assert.NoError(t, err)

	buffer.Reset()
	buffer.Write(encTransactionValidityResult)
	transactionValidityResult := primitives.DecodeTransactionValidityResult(buffer)

	assert.Equal(t, sc.Bool(true), transactionValidityResult.IsValidTransaction())
	assert.Equal(t, sc.U64(15), transactionValidityResult.AsValidTransaction().Longevity)
}

func Test_ValidateTransaction_NoUnsignedValidator(t *testing.T) {
	rt, _ := newTestRuntime(t)
	metadata := runtimeMetadata(t, rt)

	txSource := primitives.NewTransactionSourceExternal()
	blockHash := sc.BytesToFixedSequenceU8(parentHash.ToBytes())

	alice, err := ctypes.NewMultiAddressFromAccountID(signature.TestKeyringPairAlice.PublicKey)
	assert.NoError(t, err)

	amount := ctypes.NewUCompactFromUInt(constants.Dollar)

	var tests = []struct {
		callName string
		args     []any
	}{
		{
			callName: "System.remark",
			args:     []any{[]byte{}},
		},
		{
			callName: "Balances.transfer",
			args:     []any{alice, amount},
		},
		{
			callName: "Balances.set_balance",
			args:     []any{alice, amount, amount},
		},
		{
			callName: "Balances.force_transfer",
			args:     []any{alice, alice, amount},
		},
		{
			callName: "Balances.transfer_keep_alive",
			args:     []any{alice, amount},
		},
		{
			callName: "Balances.transfer_all",
			args:     []any{alice, ctypes.NewBool(false)},
		},
		{
			callName: "Balances.force_unreserve",
			args:     []any{alice, ctypes.NewU128(*big.NewInt(amount.Int64()))},
		},
	}

	for _, test := range tests {
		t.Run(test.callName, func(t *testing.T) {
			call, err := ctypes.NewCall(metadata, test.callName, test.args...)

			extrinsic := ctypes.NewExtrinsic(call)

			buffer := &bytes.Buffer{}
			txSource.Encode(buffer)

			encoder := cscale.NewEncoder(buffer)
			err = extrinsic.Encode(*encoder)
			assert.NoError(t, err)

			blockHash.Encode(buffer)

			res, err := rt.Exec("TaggedTransactionQueue_validate_transaction", buffer.Bytes())

			assert.NoError(t, err)
			assert.Equal(t,
				primitives.NewTransactionValidityResult(
					primitives.NewTransactionValidityError(
						primitives.NewUnknownTransactionNoUnsignedValidator(),
					),
				).Bytes(),
				res,
			)
		})
	}
}

func Test_ValidateTransaction_MandatoryValidation_Timestamp(t *testing.T) {
	rt, _ := newTestRuntime(t)

	idata := gossamertypes.NewInherentData()
	time := time.Now().UnixMilli()

	err := idata.SetInherent(gossamertypes.Timstap0, uint64(time))
	assert.NoError(t, err)

	ienc, err := idata.Encode()
	assert.NoError(t, err)
	inherentExt, err := rt.Exec("BlockBuilder_inherent_extrinsics", ienc)
	assert.NoError(t, err)

	txSource := primitives.NewTransactionSourceExternal()
	blockHash := sc.BytesToFixedSequenceU8(parentHash.ToBytes())

	buffer := &bytes.Buffer{}
	txSource.Encode(buffer)

	_, err = buffer.Write(inherentExt[1:])
	assert.NoError(t, err)

	blockHash.Encode(buffer)

	res, err := rt.Exec("TaggedTransactionQueue_validate_transaction", buffer.Bytes())
	assert.NoError(t, err)

	assert.NoError(t, err)
	assert.Equal(
		t,
		primitives.NewTransactionValidityResult(
			primitives.NewTransactionValidityError(
				primitives.NewInvalidTransactionMandatoryValidation(),
			),
		).Bytes(),
		res,
	)
}
