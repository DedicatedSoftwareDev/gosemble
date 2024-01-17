package module

import (
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants/testable"
	"github.com/LimeChain/gosemble/frame/testable/dispatchables"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

type TestableModule struct {
	functions map[sc.U8]primitives.Call
}

func NewTestingModule() TestableModule {
	functions := make(map[sc.U8]primitives.Call)
	functions[testable.FunctionTestIndex] = dispatchables.NewTestCall(nil)

	return TestableModule{
		functions: functions,
	}
}

func (tm TestableModule) Functions() map[sc.U8]primitives.Call {
	return tm.functions
}

func (tm TestableModule) PreDispatch(_ primitives.Call) (sc.Empty, primitives.TransactionValidityError) {
	return sc.Empty{}, nil
}

func (tm TestableModule) ValidateUnsigned(_ primitives.TransactionSource, _ primitives.Call) (primitives.ValidTransaction, primitives.TransactionValidityError) {
	return primitives.ValidTransaction{}, primitives.NewTransactionValidityError(primitives.NewUnknownTransactionNoUnsignedValidator())
}

func (tm TestableModule) Metadata() (sc.Sequence[primitives.MetadataType], primitives.MetadataModule) {
	// TODO: types
	return sc.Sequence[primitives.MetadataType]{}, primitives.MetadataModule{
		Name:      "Testable",
		Storage:   sc.Option[primitives.MetadataModuleStorage]{},
		Call:      sc.NewOption[sc.Compact](nil),
		Event:     sc.NewOption[sc.Compact](nil),
		Constants: sc.Sequence[primitives.MetadataModuleConstant]{},
		Error:     sc.NewOption[sc.Compact](nil),
		Index:     testable.ModuleIndex,
	}
}
