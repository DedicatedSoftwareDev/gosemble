package timestamp

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants"
	timestampConstants "github.com/LimeChain/gosemble/constants/timestamp"
	"github.com/LimeChain/gosemble/execution/types"
	timestamp "github.com/LimeChain/gosemble/frame/timestamp/dispatchables"
	"github.com/LimeChain/gosemble/primitives/hashing"
	"github.com/LimeChain/gosemble/primitives/log"
	"github.com/LimeChain/gosemble/primitives/storage"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

func CreateInherent(inherent primitives.InherentData) []byte {
	inherentData := inherent.Data[timestampConstants.InherentIdentifier]

	if inherentData == nil {
		log.Critical("Timestamp inherent must be provided.")
	}

	buffer := &bytes.Buffer{}
	buffer.Write(sc.SequenceU8ToBytes(inherentData))
	ts := sc.DecodeU64(buffer)
	// TODO: err if not able to parse it.
	buffer.Reset()

	timestampHash := hashing.Twox128(constants.KeyTimestamp)
	nowHash := hashing.Twox128(constants.KeyNow)

	nextTimestamp := storage.GetDecode(append(timestampHash, nowHash...), sc.DecodeU64) + timestampConstants.MinimumPeriod

	if ts > nextTimestamp {
		nextTimestamp = ts
	}

	function := timestamp.NewSetCall(sc.NewVaryingData(sc.ToCompact(uint64(nextTimestamp))))

	extrinsic := types.UncheckedExtrinsic{
		Version:  types.ExtrinsicFormatVersion,
		Function: function,
	}

	return extrinsic.Bytes()
}

func CheckInherent(args sc.VaryingData, inherent primitives.InherentData) error {
	compactTs := args[0].(sc.Compact)
	t := sc.U64(compactTs.ToBigInt().Uint64())

	inherentData := inherent.Data[timestampConstants.InherentIdentifier]

	if inherentData == nil {
		log.Critical("Timestamp inherent must be provided.")
	}

	buffer := &bytes.Buffer{}
	buffer.Write(sc.SequenceU8ToBytes(inherentData))
	ts := sc.DecodeU64(buffer)
	// TODO: err if not able to parse it.
	buffer.Reset()

	timestampHash := hashing.Twox128(constants.KeyTimestamp)
	nowHash := hashing.Twox128(constants.KeyNow)
	systemNow := storage.GetDecode(append(timestampHash, nowHash...), sc.DecodeU64)

	minimum := systemNow + timestampConstants.MinimumPeriod
	if t > ts+timestampConstants.MaxTimestampDriftMillis {
		return primitives.NewTimestampErrorTooFarInFuture()
	} else if t < minimum {
		return primitives.NewTimestampErrorTooEarly()
	}

	return nil
}
