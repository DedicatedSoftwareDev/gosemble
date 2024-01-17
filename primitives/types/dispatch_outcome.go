package types

import (
	"bytes"
	"reflect"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/log"
)

// DispatchOutcome This type specifies the outcome of dispatching a call to a module.
//
// In case of failure an error specific to the module is returned.
//
// Failure of the module call dispatching doesn't invalidate the extrinsic and it is still included
// in the block, therefore all state changes performed by the dispatched call are still persisted.
//
// For example, if the dispatching of an extrinsic involves inclusion fee payment then these
// changes are going to be preserved even if the call dispatched failed.
type DispatchOutcome sc.VaryingData //  = sc.Result[sc.Empty, DispatchError]

func NewDispatchOutcome(value sc.Encodable) DispatchOutcome {
	// None 			   = 0 - Extrinsic is valid and was submitted successfully.
	// DispatchError = 1 - Possible errors while dispatching the extrinsic.
	switch value.(type) {
	case DispatchError:
		return DispatchOutcome(sc.NewVaryingData(value))
	case sc.Empty, nil:
		return DispatchOutcome(sc.NewVaryingData(sc.Empty{}))
	default:
		log.Critical("invalid DispatchOutcome type")
	}

	panic("unreachable")
}

func (o DispatchOutcome) Encode(buffer *bytes.Buffer) {
	value := o[0]

	switch reflect.TypeOf(value) {
	case reflect.TypeOf(*new(sc.Empty)):
		sc.U8(0).Encode(buffer)
	case reflect.TypeOf(*new(DispatchError)):
		sc.U8(1).Encode(buffer)
		value.Encode(buffer)
	default:
		log.Critical("invalid DispatchOutcome type")
	}
}

func DecodeDispatchOutcome(buffer *bytes.Buffer) DispatchOutcome {
	b := sc.DecodeU8(buffer)

	switch b {
	case 0:
		return NewDispatchOutcome(sc.Empty{})
	case 1:
		value := DecodeDispatchError(buffer)
		return NewDispatchOutcome(value)
	default:
		log.Critical("invalid DispatchOutcome type")
	}

	panic("unreachable")
}

func (o DispatchOutcome) Bytes() []byte {
	return sc.EncodedBytes(o)
}
