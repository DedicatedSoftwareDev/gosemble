package extrinsic

import (
	system "github.com/LimeChain/gosemble/frame/system/extensions"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

// NewSignedPayload creates a new `SignedPayload`.
// It may fail if `additional_signed` of `Extra` is not available.
func NewSignedPayload(call primitives.Call, extra primitives.SignedExtra) (primitives.SignedPayload, primitives.TransactionValidityError) {
	additionalSigned, err := system.Extra(extra).AdditionalSigned()
	if err != nil {
		return primitives.SignedPayload{}, err
	}

	return primitives.SignedPayload{
		Call:             call,
		Extra:            extra,
		AdditionalSigned: additionalSigned,
	}, nil
}
