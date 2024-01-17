package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
)

type LastRuntimeUpgradeInfo struct {
	SpecVersion sc.Compact
	SpecName    sc.Str
}

func (lrui LastRuntimeUpgradeInfo) Encode(buffer *bytes.Buffer) {
	lrui.SpecVersion.Encode(buffer)
	lrui.SpecName.Encode(buffer)
}

func (lrui LastRuntimeUpgradeInfo) Bytes() []byte {
	buf := &bytes.Buffer{}
	lrui.Encode(buf)

	return buf.Bytes()
}

func DecodeLastRuntimeUpgradeInfo(buffer *bytes.Buffer) (value LastRuntimeUpgradeInfo) {
	if buffer.Len() <= 1 {
		return value
	}

	value.SpecVersion = sc.DecodeCompact(buffer)
	value.SpecName = sc.DecodeStr(buffer)

	return value
}
