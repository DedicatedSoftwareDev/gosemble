package main

import (
	"bytes"
	"testing"

	primitives "github.com/LimeChain/gosemble/primitives/types"

	sc "github.com/LimeChain/goscale"

	cscale "github.com/centrifuge/go-substrate-rpc-client/v4/scale"

	"github.com/centrifuge/go-substrate-rpc-client/v4/signature"
	ctypes "github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/stretchr/testify/assert"
)

func Test_TransactionPaymentApi_QueryInfo_Signed_Success(t *testing.T) {
	rt, _ := newTestRuntime(t)

	runtimeVersion, err := rt.Version()
	assert.NoError(t, err)

	metadata := runtimeMetadata(t, rt)

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

	buffer := &bytes.Buffer{}
	encoder := cscale.NewEncoder(buffer)
	err = extrinsic.Encode(*encoder)
	assert.NoError(t, err)

	sc.U32(buffer.Len()).Encode(buffer)

	bytesRuntimeDispatchInfo, err := rt.Exec("TransactionPaymentApi_query_info", buffer.Bytes())
	assert.NoError(t, err)

	buffer.Reset()
	buffer.Write(bytesRuntimeDispatchInfo)

	rdi := primitives.DecodeRuntimeDispatchInfo(buffer)

	expectedRdi := primitives.RuntimeDispatchInfo{
		Weight:     primitives.WeightFromParts(2_091_000, 0),
		Class:      primitives.NewDispatchClassNormal(),
		PartialFee: sc.NewU128FromUint64(110_536_107),
	}

	assert.Equal(t, expectedRdi, rdi)
}

func Test_TransactionPaymentApi_QueryInfo_Unsigned_Success(t *testing.T) {
	rt, _ := newTestRuntime(t)

	metadata := runtimeMetadata(t, rt)

	call, err := ctypes.NewCall(metadata, "System.remark", []byte{})
	assert.NoError(t, err)

	extrinsic := ctypes.NewExtrinsic(call)

	buffer := &bytes.Buffer{}
	encoder := cscale.NewEncoder(buffer)
	err = extrinsic.Encode(*encoder)
	assert.NoError(t, err)

	sc.U32(buffer.Len()).Encode(buffer)

	bytesRuntimeDispatchInfo, err := rt.Exec("TransactionPaymentApi_query_info", buffer.Bytes())
	assert.NoError(t, err)

	buffer.Reset()
	buffer.Write(bytesRuntimeDispatchInfo)

	rdi := primitives.DecodeRuntimeDispatchInfo(buffer)

	expectedRdi := primitives.RuntimeDispatchInfo{
		Weight:     primitives.WeightFromParts(2_091_000, 0),
		Class:      primitives.NewDispatchClassNormal(),
		PartialFee: sc.NewU128FromUint64(0),
	}

	assert.Equal(t, expectedRdi, rdi)
}

func Test_TransactionPaymentApi_QueryFeeDetails_Signed_Success(t *testing.T) {
	rt, _ := newTestRuntime(t)

	runtimeVersion, err := rt.Version()
	assert.NoError(t, err)

	metadata := runtimeMetadata(t, rt)

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

	buffer := &bytes.Buffer{}
	encoder := cscale.NewEncoder(buffer)
	err = extrinsic.Encode(*encoder)
	assert.NoError(t, err)

	sc.U32(buffer.Len()).Encode(buffer)

	bytesFeeDetails, err := rt.Exec("TransactionPaymentApi_query_fee_details", buffer.Bytes())
	assert.NoError(t, err)

	buffer.Reset()
	buffer.Write(bytesFeeDetails)

	fd := primitives.DecodeFeeDetails(buffer)

	expectedFd := primitives.FeeDetails{
		InclusionFee: sc.NewOption[primitives.InclusionFee](
			primitives.NewInclusionFee(
				sc.NewU128FromUint64(110_536_000),
				sc.NewU128FromUint64(107),
				sc.NewU128FromUint64(0),
			)),
	}

	assert.Equal(t, expectedFd, fd)
}

func Test_TransactionPaymentApi_QueryFeeDetails_Unsigned_Success(t *testing.T) {
	rt, _ := newTestRuntime(t)
	metadata := runtimeMetadata(t, rt)

	call, err := ctypes.NewCall(metadata, "System.remark", []byte{})
	assert.NoError(t, err)

	extrinsic := ctypes.NewExtrinsic(call)

	buffer := &bytes.Buffer{}
	encoder := cscale.NewEncoder(buffer)
	err = extrinsic.Encode(*encoder)
	assert.NoError(t, err)

	sc.U32(buffer.Len()).Encode(buffer)

	bytesFeeDetails, err := rt.Exec("TransactionPaymentApi_query_fee_details", buffer.Bytes())
	assert.NoError(t, err)

	buffer.Reset()
	buffer.Write(bytesFeeDetails)

	fd := primitives.DecodeFeeDetails(buffer)

	expectedFd := primitives.FeeDetails{
		InclusionFee: sc.NewOption[primitives.InclusionFee](nil),
	}

	assert.Equal(t, expectedFd, fd)
}

func Test_TransactionPaymentCallApi_QueryCallInfo_Success(t *testing.T) {
	rt, _ := newTestRuntime(t)
	metadata := runtimeMetadata(t, rt)

	call, err := ctypes.NewCall(metadata, "System.remark", []byte{})
	assert.NoError(t, err)

	buffer := &bytes.Buffer{}
	encoder := cscale.NewEncoder(buffer)
	err = call.CallIndex.Encode(*encoder)
	assert.NoError(t, err)

	err = call.Args.Encode(*encoder)
	assert.NoError(t, err)

	sc.U32(buffer.Len()).Encode(buffer)

	bytesRuntimeDispatchInfo, err := rt.Exec("TransactionPaymentCallApi_query_call_info", buffer.Bytes())
	assert.NoError(t, err)

	buffer.Reset()
	buffer.Write(bytesRuntimeDispatchInfo)

	rdi := primitives.DecodeRuntimeDispatchInfo(buffer)

	expectedRdi := primitives.RuntimeDispatchInfo{
		Weight:     primitives.WeightFromParts(2_091_000, 0),
		Class:      primitives.NewDispatchClassNormal(),
		PartialFee: sc.NewU128FromUint64(110_536_003),
	}

	assert.Equal(t, expectedRdi, rdi)
}

func Test_TransactionPaymentCallApi_QueryCallFeeDetails_Success(t *testing.T) {
	rt, _ := newTestRuntime(t)
	metadata := runtimeMetadata(t, rt)

	call, err := ctypes.NewCall(metadata, "System.remark", []byte{})
	assert.NoError(t, err)

	buffer := &bytes.Buffer{}
	encoder := cscale.NewEncoder(buffer)
	err = call.CallIndex.Encode(*encoder)
	assert.NoError(t, err)

	err = call.Args.Encode(*encoder)
	assert.NoError(t, err)

	sc.U32(buffer.Len()).Encode(buffer)

	bytesFeeDetails, err := rt.Exec("TransactionPaymentCallApi_query_call_fee_details", buffer.Bytes())
	assert.NoError(t, err)

	buffer.Reset()
	buffer.Write(bytesFeeDetails)

	fd := primitives.DecodeFeeDetails(buffer)

	expectedFd := primitives.FeeDetails{
		InclusionFee: sc.NewOption[primitives.InclusionFee](
			primitives.NewInclusionFee(
				sc.NewU128FromUint64(110_536_000),
				sc.NewU128FromUint64(3),
				sc.NewU128FromUint64(0),
			)),
	}

	assert.Equal(t, expectedFd, fd)
}
