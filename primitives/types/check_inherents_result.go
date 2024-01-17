package types

import (
	"bytes"
	"fmt"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/log"
)

const (
	InherentErrorInherentDataExists sc.U8 = iota
	InherentErrorDecodingFailed
	InherentErrorFatalErrorReported
	InherentErrorApplication
)

const (
	errDecodeInherentData = "failed to decode InherentData"
)

type InherentError struct {
	sc.VaryingData
}

func NewInherentErrorInherentDataExists(inherentIdentifier sc.Sequence[sc.U8]) InherentError {
	return InherentError{sc.NewVaryingData(InherentErrorInherentDataExists, inherentIdentifier)}
}

func NewInherentErrorDecodingFailed(inherentIdentifier sc.Sequence[sc.U8]) InherentError {
	return InherentError{sc.NewVaryingData(InherentErrorDecodingFailed, inherentIdentifier)}
}

func NewInherentErrorFatalErrorReported() InherentError {
	return InherentError{sc.NewVaryingData(InherentErrorFatalErrorReported)}
}

func NewInherentErrorApplication() InherentError {
	// TODO: encode additional value
	return InherentError{sc.NewVaryingData(InherentErrorApplication)}
}

func (ie InherentError) IsFatal() sc.Bool {
	switch ie.VaryingData[0] {
	case InherentErrorFatalErrorReported:
		return true
	default:
		return false
	}
}

func (ie InherentError) Error() string {
	switch ie.VaryingData[0] {
	case InherentErrorInherentDataExists:
		return fmt.Sprintf("Inherent data already exists for identifier: [%v]", ie.VaryingData[1])
	case InherentErrorDecodingFailed:
		return fmt.Sprintf("Failed to decode inherent data for identifier: [%v]", ie.VaryingData[1])
	case InherentErrorFatalErrorReported:
		return "There was already a fatal error reported and no other errors are allowed"
	case InherentErrorApplication:
		return "Inherent error application"
	default:
		log.Critical("invalid inherent error")
	}

	panic("unreachable")
}

type CheckInherentsResult struct {
	Okay       sc.Bool
	FatalError sc.Bool
	Errors     InherentData
}

func NewCheckInherentsResult() CheckInherentsResult {
	return CheckInherentsResult{
		Okay:       true,
		FatalError: false,
		Errors:     *NewInherentData(),
	}
}

func (cir CheckInherentsResult) Encode(buffer *bytes.Buffer) {
	cir.Okay.Encode(buffer)
	cir.FatalError.Encode(buffer)
	cir.Errors.Encode(buffer)
}

func (cir CheckInherentsResult) PutError(inherentIdentifier [8]byte, error IsFatalError) error {
	if cir.FatalError {
		return NewInherentErrorFatalErrorReported()
	}

	if error.IsFatal() {
		cir.Errors.Clear()
	}

	err := cir.Errors.Put(inherentIdentifier, error)
	if err != nil {
		return err
	}

	cir.Okay = false
	cir.FatalError = error.IsFatal()

	return nil
}

func DecodeCheckInherentsResult(buffer *bytes.Buffer) CheckInherentsResult {
	okay := sc.DecodeBool(buffer)
	fatalError := sc.DecodeBool(buffer)
	errors, err := DecodeInherentData(buffer)
	if err != nil {
		log.Critical(errDecodeInherentData)
	}

	return CheckInherentsResult{
		Okay:       okay,
		FatalError: fatalError,
		Errors:     *errors,
	}
}
