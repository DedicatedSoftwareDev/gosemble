package extrinsic

import (
	sc "github.com/LimeChain/goscale"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

type Applyable interface {
	Apply(validator UnsignedValidator, info *primitives.DispatchInfo, length sc.Compact) (ok primitives.DispatchResultWithPostInfo[primitives.PostDispatchInfo], err primitives.TransactionValidityError)
}

type Validatable interface {
	Validate(validator UnsignedValidator, source primitives.TransactionSource, info *primitives.DispatchInfo, length sc.Compact) (ok primitives.ValidTransaction, err primitives.TransactionValidityError)
}

// UnsignedValidator provides validation for unsigned extrinsics.
//
// This trait provides two functions [`pre_dispatch`](Self::pre_dispatch) and
// [`validate_unsigned`](Self::validate_unsigned). The [`pre_dispatch`](Self::pre_dispatch)
// function is called right before dispatching the call wrapped by an unsigned extrinsic. The
// [`validate_unsigned`](Self::validate_unsigned) function is mainly being used in the context of
// the transaction pool to check the validity of the call wrapped by an unsigned extrinsic.
type UnsignedValidator interface {
	// PreDispatch validates the call right before dispatch.
	//
	// This method should be used to prevent transactions already in the pool
	// (i.e. passing [`validate_unsigned`](Self::validate_unsigned)) from being included in blocks
	// in case they became invalid since being added to the pool.
	//
	// By default it's a good idea to call [`validate_unsigned`](Self::validate_unsigned) from
	// within this function again to make sure we never include an invalid transaction. Otherwise
	// the implementation of the call or this method will need to provide proper validation to
	// ensure that the transaction is valid.
	//
	// Changes made to storage *WILL* be persisted if the call returns `Ok`.
	PreDispatch(call *primitives.Call) (ok sc.Empty, err primitives.TransactionValidityError)

	// ValidateUnsigned returns the validity of the call
	//
	// This method has no side-effects. It merely checks whether the call would be rejected
	// by the runtime in an unsigned extrinsic.
	//
	// The validity checks should be as lightweight as possible because every node will execute
	// this code before the unsigned extrinsic enters the transaction pool and also periodically
	// afterwards to ensure the validity. To prevent dos-ing a network with unsigned
	// extrinsics, these validity checks should include some checks around uniqueness, for example,
	// like checking that the unsigned extrinsic was send by an authority in the active set.
	//
	// Changes made to storage should be discarded by caller.
	ValidateUnsigned(source primitives.TransactionSource, call *primitives.Call) (ok primitives.ValidTransaction, err primitives.TransactionValidityError)
}
