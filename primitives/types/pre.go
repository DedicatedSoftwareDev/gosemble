package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
)

// Pre is the type that encodes information that can be passed from pre_dispatch to post-dispatch.
type Pre struct {
	Tip       Balance
	Who       Address32
	Imbalance sc.Option[Balance]
}

func (p Pre) Encode(buffer *bytes.Buffer) {
	p.Tip.Encode(buffer)
	p.Who.Encode(buffer)
	p.Imbalance.Encode(buffer)
}

func (p Pre) Bytes() []byte {
	return sc.EncodedBytes(p)
}
