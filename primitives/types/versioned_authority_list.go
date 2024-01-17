package types

import (
	"bytes"
	sc "github.com/LimeChain/goscale"
)

type VersionedAuthorityList struct {
	Version       sc.U8
	AuthorityList sc.Sequence[Authority]
}

func (val VersionedAuthorityList) Encode(buffer *bytes.Buffer) {
	val.Version.Encode(buffer)
	val.AuthorityList.Encode(buffer)
}

func DecodeVersionedAuthorityList(buffer *bytes.Buffer) VersionedAuthorityList {
	return VersionedAuthorityList{
		Version:       sc.DecodeU8(buffer),
		AuthorityList: sc.DecodeSequenceWith(buffer, DecodeAuthority),
	}
}

func (val VersionedAuthorityList) Bytes() []byte {
	return sc.EncodedBytes(val)
}
