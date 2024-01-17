package types

// TODO:
//// Means by which a transaction may be extended. This type embodies both the data and the logic
//// that should be additionally associated with the transaction. It should be plain old data.
//type SignedExtension interface {
//	sc.Encodable
//
//	// Unique identifier of this signed extension.
//	//
//	// This will be exposed in the metadata to identify the signed extension used
//	// in an extrinsic.
//	// const IDENTIFIER: &'static str;
//
//	// The type which encodes the sender identity.
//	// type AccountId;
//
//	// The type which encodes the call to be dispatched.
//	// type Call: Dispatchable
//
//	// Any additional data that will go into the signed payload. This may be created dynamically
//	// from the transaction using the `additional_signed` function.
//	// type AdditionalSigned: Encode + TypeInfo
//
//	// The type that encodes information that can be passed from pre_dispatch to post-dispatch.
//	// type Pre
//
//	// Construct any additional data that should be in the signed payload of the transaction. Can
//	// also perform any pre-signature-verification checks and return an error if needed.
//	AdditionalSigned() (ok AdditionalSigned, err TransactionValidityError)
//
//	// Validate a signed transaction for the transaction queue.
//	//
//	// This function can be called frequently by the transaction queue,
//	// to obtain transaction validity against current state.
//	// It should perform all checks that determine a valid transaction,
//	// that can pay for its execution and quickly eliminate ones
//	// that are stale or incorrect.
//	//
//	// Make sure to perform the same checks in `pre_dispatch` function.
//	Validate(_who *Address32, _call *Call, _info *DispatchInfo, _length sc.Compact) (ok ValidTransaction, err TransactionValidityError)
//
//	// Do any pre-flight stuff for a signed transaction.
//	//
//	// Make sure to perform the same checks as in [`Self::validate`].
//	PreDispatch(e SignedExtra, who *Address32, call *Call, info *DispatchInfo, length sc.Compact) (ok Pre, err TransactionValidityError)
//
//	// Validate an unsigned transaction for the transaction queue.
//	//
//	// This function can be called frequently by the transaction queue
//	// to obtain transaction validity against current state.
//	// It should perform all checks that determine a valid unsigned transaction,
//	// and quickly eliminate ones that are stale or incorrect.
//	//
//	// Make sure to perform the same checks in `pre_dispatch_unsigned` function.
//	ValidateUnsigned(_call *Call, _info *DispatchInfo, _length sc.Compact) (ok ValidTransaction, err TransactionValidityError)
//
//	// Do any pre-flight stuff for a unsigned transaction.
//	//
//	// Note this function by default delegates to `validate_unsigned`, so that
//	// all checks performed for the transaction queue are also performed during
//	// the dispatch phase (applying the extrinsic).
//	//
//	// If you ever override this function, you need to make sure to always
//	// perform the same validation as in `validate_unsigned`.
//	PreDispatchUnsigned(call *Call, info *DispatchInfo, length sc.Compact) (ok Pre, err TransactionValidityError)
//
//	// Do any post-flight stuff for an extrinsic.
//	//
//	// If the transaction is signed, then `_pre` will contain the output of `pre_dispatch`,
//	// and `None` otherwise.
//	//
//	// This gets given the `DispatchResult` `_result` from the extrinsic and can, if desired,
//	// introduce a `TransactionValidityError`, causing the block to become invalid for including
//	// it.
//	//
//	// WARNING: It is dangerous to return an error here. To do so will fundamentally invalidate the
//	// transaction and any block that it is included in, causing the block author to not be
//	// compensated for their work in validating the transaction or producing the block so far.
//	//
//	// It can only be used safely when you *know* that the extrinsic is one that can only be
//	// introduced by the current block author; generally this implies that it is an inherent and
//	// will come from either an offchain-worker or via `InherentData`.(
//	PostDispatch(pre sc.Option[Pre], info *DispatchInfo, postInfo *PostDispatchInfo, length sc.Compact, result *DispatchResult) (ok Pre, err TransactionValidityError)
//}
