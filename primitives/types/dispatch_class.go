package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/log"
)

const (
	// DispatchClassNormal A normal dispatch.
	DispatchClassNormal sc.U8 = iota

	// DispatchClassOperational An operational dispatch.
	DispatchClassOperational

	// DispatchClassMandatory A mandatory dispatch. These kinds of dispatch are always included regardless of their
	// weight, therefore it is critical that they are separately validated to ensure that a
	// malicious validator cannot craft a valid but impossibly heavy block. Usually this just
	// means ensuring that the extrinsic can only be included once and that it is always very
	// light.
	//
	// Do *NOT* use it for extrinsics that can be heavy.
	//
	// The only real use case for this is inherent extrinsics that are required to execute in a
	// block for the block to be valid, and it solves the issue in the case that the block
	// initialization is sufficiently heavy to mean that those inherents do not fit into the
	// block. Essentially, we assume that in these exceptional circumstances, it is better to
	// allow an overweight block to be created than to not allow any block at all to be created.
	DispatchClassMandatory
)

// A generalized group of dispatch types.
type DispatchClass struct {
	sc.VaryingData
}

func NewDispatchClassNormal() DispatchClass {
	return DispatchClass{sc.NewVaryingData(DispatchClassNormal)}
}

func NewDispatchClassOperational() DispatchClass {
	return DispatchClass{sc.NewVaryingData(DispatchClassOperational)}
}

func NewDispatchClassMandatory() DispatchClass {
	return DispatchClass{sc.NewVaryingData(DispatchClassMandatory)}
}

func DecodeDispatchClass(buffer *bytes.Buffer) DispatchClass {
	b := sc.DecodeU8(buffer)

	switch b {
	case DispatchClassNormal:
		return NewDispatchClassNormal()
	case DispatchClassOperational:
		return NewDispatchClassOperational()
	case DispatchClassMandatory:
		return NewDispatchClassMandatory()
	default:
		log.Critical("invalid DispatchClass type")
	}

	panic("unreachable")
}

func (dc DispatchClass) Is(value sc.U8) sc.Bool {
	// TODO: type safety
	switch value {
	case DispatchClassNormal, DispatchClassOperational, DispatchClassMandatory:
		return dc.VaryingData[0] == value
	default:
		log.Critical("invalid DispatchClass value")
	}

	panic("unreachable")
}

// Returns an array containing all dispatch classes.
func DispatchClassAll() []DispatchClass {
	return []DispatchClass{NewDispatchClassNormal(), NewDispatchClassOperational(), NewDispatchClassMandatory()}
}

// A struct holding value for each `DispatchClass`.
type PerDispatchClass[T sc.Encodable] struct {
	// Value for `Normal` extrinsics.
	Normal T
	// Value for `Operational` extrinsics.
	Operational T
	// Value for `Mandatory` extrinsics.
	Mandatory T
}

func (pdc PerDispatchClass[T]) Encode(buffer *bytes.Buffer) {
	pdc.Normal.Encode(buffer)
	pdc.Operational.Encode(buffer)
	pdc.Mandatory.Encode(buffer)
}

func DecodePerDispatchClass[T sc.Encodable](buffer *bytes.Buffer, decodeFunc func(buffer *bytes.Buffer) T) PerDispatchClass[T] {
	return PerDispatchClass[T]{
		Normal:      decodeFunc(buffer),
		Operational: decodeFunc(buffer),
		Mandatory:   decodeFunc(buffer),
	}
}

func (pdc PerDispatchClass[T]) Bytes() []byte {
	return sc.EncodedBytes(pdc)
}

// Get current value for given class.
func (pdc *PerDispatchClass[T]) Get(class DispatchClass) *T {
	switch class.VaryingData[0] {
	case DispatchClassNormal:
		return &pdc.Normal
	case DispatchClassOperational:
		return &pdc.Operational
	case DispatchClassMandatory:
		return &pdc.Mandatory
	default:
		log.Critical("invalid DispatchClass type")
	}

	panic("unreachable")
}

// An object to track the currently used extrinsic weight in a block.
type ConsumedWeight PerDispatchClass[Weight]

func (cw ConsumedWeight) Encode(buffer *bytes.Buffer) {
	PerDispatchClass[Weight](cw).Encode(buffer)
}

func DecodeConsumedWeight(buffer *bytes.Buffer) ConsumedWeight {
	return ConsumedWeight{
		Normal:      DecodeWeight(buffer),
		Operational: DecodeWeight(buffer),
		Mandatory:   DecodeWeight(buffer),
	}
}

func (cw ConsumedWeight) Bytes() []byte {
	return sc.EncodedBytes(cw)
}

// Get current value for given class.
func (cw *ConsumedWeight) Get(class DispatchClass) *Weight {
	switch class.VaryingData[0] {
	case DispatchClassNormal:
		return &cw.Normal
	case DispatchClassOperational:
		return &cw.Operational
	case DispatchClassMandatory:
		return &cw.Mandatory
	default:
		log.Critical("invalid DispatchClass type")
	}

	panic("unreachable")
}

// Returns the total weight consumed by all extrinsics in the block.
//
// Saturates on overflow.
func (cw ConsumedWeight) Total() Weight {
	sum := WeightZero()
	for _, class := range []DispatchClass{NewDispatchClassNormal(), NewDispatchClassOperational(), NewDispatchClassMandatory()} {
		sum = sum.SaturatingAdd(*cw.Get(class))
	}
	return sum
}

// SaturatingAdd Increase the weight of the given class. Saturates at the numeric bounds.
func (cw *ConsumedWeight) SaturatingAdd(weight Weight, class DispatchClass) {
	weightForClass := cw.Get(class)
	weightForClass.RefTime = weightForClass.RefTime.SaturatingAdd(weight.RefTime)
	weightForClass.ProofSize = weightForClass.ProofSize.SaturatingAdd(weight.ProofSize)
}

// Accrue Increase the weight of the given class. Saturates at the numeric bounds.
func (cw *ConsumedWeight) Accrue(weight Weight, class DispatchClass) {
	weightForClass := cw.Get(class)
	weightForClass.SaturatingAccrue(weight)
}

// CheckedAccrue Try to increase the weight of the given class. Saturates at the numeric bounds.
func (cw *ConsumedWeight) CheckedAccrue(weight Weight, class DispatchClass) (ok sc.Empty, err error) {
	weightForClass := cw.Get(class)
	refTime, err := weightForClass.RefTime.CheckedAdd(weight.RefTime)
	if err != nil {
		return ok, err
	}
	weightForClass.RefTime = refTime

	proofSize, err := weightForClass.ProofSize.CheckedAdd(weight.ProofSize)
	if err != nil {
		return ok, err
	}
	weightForClass.ProofSize = proofSize

	return ok, err
}

// Reduce the weight of the given class. Saturates at the numeric bounds.
func (cw *ConsumedWeight) Reduce(weight Weight, class DispatchClass) {
	weightForClass := cw.Get(class)
	weightForClass.SaturatingReduce(weight)
}
