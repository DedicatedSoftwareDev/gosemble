package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

// CheckedExtrinsic is the definition of something that the external world might want to say; its
// existence implies that it has been checked and is good, particularly with
// regards to the signature.
//
// TODO: make it generic
// generic::CheckedExtrinsic<AccountId, RuntimeCall, SignedExtra>;
type CheckedExtrinsic struct {
	Version sc.U8

	// Who this purports to be from and the number of extrinsics have come before
	// from the same signer, if anyone (note this is not a signature).
	Signed   sc.Option[primitives.AccountIdExtra]
	Function primitives.Call
}

func (cxt CheckedExtrinsic) Encode(buffer *bytes.Buffer) {
	cxt.Version.Encode(buffer)
	cxt.Signed.Encode(buffer)
	cxt.Function.Encode(buffer)
}

func (cxt CheckedExtrinsic) Bytes() []byte {
	return sc.EncodedBytes(cxt)
}
