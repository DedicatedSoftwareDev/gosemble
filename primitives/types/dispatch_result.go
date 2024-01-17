package types

import (
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/log"
)

type DispatchResult sc.VaryingData

func NewDispatchResult(value sc.Encodable) DispatchResult {
	switch value.(type) {
	case DispatchError, DispatchErrorWithPostInfo[PostDispatchInfo]:
		return DispatchResult(sc.NewVaryingData(value))
	case sc.Empty, nil:
		return DispatchResult(sc.NewVaryingData(sc.Empty{}))
	default:
		log.Critical("invalid DispatchResult type")
	}

	panic("unreachable")
}

type DispatchResultWithPostInfo[T sc.Encodable] struct {
	HasError sc.Bool
	Ok       T
	Err      DispatchErrorWithPostInfo[T]
}
