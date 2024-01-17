package extrinsic

import (
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/execution/types"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

type Unchecked types.UncheckedExtrinsic

func (uxt Unchecked) Check(lookup primitives.AccountIdLookup) (ok types.CheckedExtrinsic, err primitives.TransactionValidityError) {
	switch uxt.Signature.HasValue {
	case true:
		signer, signature, extra := uxt.Signature.Value.Signer, uxt.Signature.Value.Signature, uxt.Signature.Value.Extra

		signedAddress, err := lookup.Lookup(signer)
		if err != nil {
			return ok, err
		}

		rawPayload, err := NewSignedPayload(uxt.Function, extra)
		if err != nil {
			return ok, err
		}

		if !signature.Verify(rawPayload.UsingEncoded(), signedAddress) {
			err := primitives.NewTransactionValidityError(primitives.NewInvalidTransactionBadProof())
			return ok, err
		}

		function, extra, _ := rawPayload.Call, rawPayload.Extra, rawPayload.AdditionalSigned

		ok = types.CheckedExtrinsic{
			Signed:   sc.NewOption[primitives.AccountIdExtra](primitives.AccountIdExtra{Address32: signedAddress, SignedExtra: extra}),
			Function: function,
		}
	case false:
		ok = types.CheckedExtrinsic{
			Signed:   sc.NewOption[primitives.AccountIdExtra](nil),
			Function: uxt.Function,
		}
	}

	return ok, err
}
