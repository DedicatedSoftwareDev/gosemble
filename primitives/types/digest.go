package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
)

const (
	DigestTypeConsensusMessage           = 4
	DigestTypeSeal                       = 5
	DigestTypePreRuntime                 = 6
	DigestTypeRuntimeEnvironmentUpgraded = 8
)

type Digest = sc.Dictionary[sc.U8, sc.FixedSequence[DigestItem]]

func DecodeDigest(buffer *bytes.Buffer) Digest {
	compactSize := sc.DecodeCompact(buffer)
	size := int(compactSize.ToBigInt().Int64())

	decoder := sc.Decoder{Reader: buffer}

	result := Digest{}
	for i := 0; i < size; i++ {
		digestType := decoder.DecodeByte()

		switch digestType {
		case DigestTypeConsensusMessage:
			consensusDigest := DecodeDigestItem(buffer)
			result[DigestTypeConsensusMessage] = append(result[DigestTypeConsensusMessage], consensusDigest)
		case DigestTypeSeal:
			seal := DecodeDigestItem(buffer)
			result[DigestTypeSeal] = append(result[DigestTypeSeal], seal)
		case DigestTypePreRuntime:
			preRuntimeDigest := DecodeDigestItem(buffer)
			result[DigestTypePreRuntime] = append(result[DigestTypePreRuntime], preRuntimeDigest)
		case DigestTypeRuntimeEnvironmentUpgraded:
			sc.DecodeU8(buffer)
			// TODO:
		}
	}

	return result
}
