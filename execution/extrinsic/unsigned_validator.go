package extrinsic

import (
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/config"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

type UnsignedValidatorForChecked struct{}

// PreDispatch validates the dispatch call before execution.
// Inherent call is accepted for being dispatched
func (v UnsignedValidatorForChecked) PreDispatch(call *primitives.Call) (sc.Empty, primitives.TransactionValidityError) {
	module, ok := config.Modules[(*call).ModuleIndex()]
	if !ok {
		return sc.Empty{}, nil
	}

	return module.PreDispatch(*call)
}

// ValidateUnsigned returns the validity of the dispatch call.
// Inherent call is not validated as unsigned
func (v UnsignedValidatorForChecked) ValidateUnsigned(_source primitives.TransactionSource, call *primitives.Call) (primitives.ValidTransaction, primitives.TransactionValidityError) {
	module, ok := config.Modules[(*call).ModuleIndex()]
	if !ok {
		return primitives.ValidTransaction{}, primitives.NewTransactionValidityError(primitives.NewUnknownTransactionNoUnsignedValidator())
	}

	return module.ValidateUnsigned(_source, *call)
}
