package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
)

// TODO: Extend for different types (ecdsa, ed25519, sr25519)
type PublicKey = sc.FixedSequence[sc.U8]

func DecodePublicKey(buffer *bytes.Buffer) PublicKey {
	return sc.DecodeFixedSequence[sc.U8](32, buffer)
}
