package dispatchables

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants"
	"github.com/LimeChain/gosemble/constants/timestamp"
	"github.com/LimeChain/gosemble/frame/aura"
	"github.com/LimeChain/gosemble/primitives/hashing"
	"github.com/LimeChain/gosemble/primitives/log"
	"github.com/LimeChain/gosemble/primitives/storage"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

type SetCall struct {
	primitives.Callable
}

func NewSetCall(args sc.VaryingData) SetCall {
	call := SetCall{
		Callable: primitives.Callable{
			ModuleId:   timestamp.ModuleIndex,
			FunctionId: timestamp.FunctionSetIndex,
		},
	}

	if len(args) != 0 {
		call.Arguments = args
	}

	return call
}

func (c SetCall) DecodeArgs(buffer *bytes.Buffer) primitives.Call {
	c.Arguments = sc.NewVaryingData(sc.DecodeCompact(buffer))
	return c
}

func (c SetCall) Encode(buffer *bytes.Buffer) {
	c.Callable.Encode(buffer)
}

func (c SetCall) Bytes() []byte {
	return c.Callable.Bytes()
}

func (c SetCall) ModuleIndex() sc.U8 {
	return c.Callable.ModuleIndex()
}

func (c SetCall) FunctionIndex() sc.U8 {
	return c.Callable.FunctionIndex()
}

func (c SetCall) Args() sc.VaryingData {
	return c.Callable.Args()
}

func (_ SetCall) BaseWeight(b ...any) primitives.Weight {
	// Storage: Timestamp Now (r:1 w:1)
	// Proof: Timestamp Now (max_values: Some(1), max_size: Some(8), added: 503, mode: MaxEncodedLen)
	// Storage: Babe CurrentSlot (r:1 w:0)
	// Proof: Babe CurrentSlot (max_values: Some(1), max_size: Some(8), added: 503, mode: MaxEncodedLen)
	// TODO: Consensus algorithm affects weight values.
	// Proof Size summary in bytes:
	//  Measured:  `312`
	//  Estimated: `1006`
	// Minimum execution time: 9_106 nanoseconds.
	r := constants.DbWeight.Reads(2)
	w := constants.DbWeight.Writes(1)
	return primitives.WeightFromParts(9_258_000, 1006).SaturatingAdd(r).SaturatingAdd(w)
}

func (_ SetCall) IsInherent() bool {
	return true
}

func (_ SetCall) WeightInfo(baseWeight primitives.Weight) primitives.Weight {
	return primitives.WeightFromParts(baseWeight.RefTime, 0)
}

func (_ SetCall) ClassifyDispatch(baseWeight primitives.Weight) primitives.DispatchClass {
	return primitives.NewDispatchClassMandatory()
}

func (_ SetCall) PaysFee(baseWeight primitives.Weight) primitives.Pays {
	return primitives.NewPaysYes()
}

func (_ SetCall) Dispatch(origin primitives.RuntimeOrigin, args sc.VaryingData) primitives.DispatchResultWithPostInfo[primitives.PostDispatchInfo] {
	compactTs := args[0].(sc.Compact)
	return set(origin, sc.U64(compactTs.ToBigInt().Uint64()))
}

// set sets the current time.
//
// This call should be invoked exactly once per block. It will panic at the finalization
// phase, if this call hasn't been invoked by that time.
//
// The timestamp should be greater than the previous one by the amount specified by
// `MinimumPeriod`.
//
// The dispatch origin for this call must be `Inherent`.
//
// ## Complexity
//   - `O(1)` (Note that implementations of `OnTimestampSet` must also be `O(1)`)
//   - 1 storage read and 1 storage mutation (codec `O(1)`). (because of `DidUpdate::take` in
//     `on_finalize`)
//   - 1 event handler `on_timestamp_set`. Must be `O(1)`.
func set(origin primitives.RuntimeOrigin, now sc.U64) primitives.DispatchResultWithPostInfo[primitives.PostDispatchInfo] {
	if !origin.IsNoneOrigin() {
		return primitives.DispatchResultWithPostInfo[primitives.PostDispatchInfo]{
			HasError: true,
			Err: primitives.DispatchErrorWithPostInfo[primitives.PostDispatchInfo]{
				Error: primitives.NewDispatchErrorBadOrigin(),
			},
		}
	}

	timestampHash := hashing.Twox128(constants.KeyTimestamp)
	didUpdateHash := hashing.Twox128(constants.KeyDidUpdate)

	didUpdate := storage.Exists(append(timestampHash, didUpdateHash...))

	if didUpdate == 1 {
		log.Critical("Timestamp must be updated only once in the block")
	}

	nowHash := hashing.Twox128(constants.KeyNow)
	previousTimestamp := storage.GetDecode(append(timestampHash, nowHash...), sc.DecodeU64)

	if !(previousTimestamp == 0 || now >= previousTimestamp+timestamp.MinimumPeriod) {
		log.Critical("Timestamp must increment by at least <MinimumPeriod> between sequential blocks")
	}

	storage.Set(append(timestampHash, nowHash...), now.Bytes())
	storage.Set(append(timestampHash, didUpdateHash...), sc.Bool(true).Bytes())

	// TODO: Every consensus that uses the timestamp must implement
	// <T::OnTimestampSet as OnTimestampSet<_>>::on_timestamp_set(now)

	// TODO:
	// timestamp module should not depend on the aura module
	aura.OnTimestampSet(now)

	return primitives.DispatchResultWithPostInfo[primitives.PostDispatchInfo]{
		HasError: false,
		Ok:       primitives.PostDispatchInfo{},
	}
}
