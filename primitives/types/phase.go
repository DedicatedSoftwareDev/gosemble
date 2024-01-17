package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/log"
)

const (
	// PhaseApplyExtrinsic Applying an extrinsic.
	PhaseApplyExtrinsic sc.U8 = iota

	// PhaseFinalization Finalizing the block.
	PhaseFinalization

	// PhaseInitialization Initializing the block.
	PhaseInitialization
)

type ExtrinsicPhase = sc.VaryingData

func NewExtrinsicPhaseApply(index sc.U32) ExtrinsicPhase {
	return sc.NewVaryingData(PhaseApplyExtrinsic, index)
}

func NewExtrinsicPhaseFinalization() ExtrinsicPhase {
	return sc.NewVaryingData(PhaseFinalization)
}

func NewExtrinsicPhaseInitialization() ExtrinsicPhase {
	return sc.NewVaryingData(PhaseInitialization)
}

func DecodeExtrinsicPhase(buffer *bytes.Buffer) ExtrinsicPhase {
	b := sc.DecodeU8(buffer)

	switch b {
	case PhaseApplyExtrinsic:
		index := sc.DecodeU32(buffer)
		return NewExtrinsicPhaseApply(index)
	case PhaseFinalization:
		return NewExtrinsicPhaseFinalization()
	case PhaseInitialization:
		return NewExtrinsicPhaseInitialization()
	default:
		log.Critical("invalid Phase type")
	}

	panic("unreachable")
}
