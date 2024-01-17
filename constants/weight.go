package constants

import (
	"math"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/types"
)

// TODO: needs to be benchmarked

const FiveMbPerBlockPerExtrinsic sc.U32 = 5 * 1024 * 1024
const WeightRefTimePerSecond sc.U64 = 1_000_000_000_000
const WeightRefTimePerNanos sc.U64 = 1_000

// We assume that ~10% of the block weight is consumed by `on_initialize` handlers.
// This is used to limit the maximal weight of a single extrinsic.
var AverageOnInitializeRatio types.Perbill = types.Perbill{Percentage: 10}

// We allow `Normal` extrinsics to fill up the block up to 75%, the rest can be used
// by  Operational  extrinsics.
var NormalDispatchRatio types.Perbill = types.Perbill{Percentage: 75}

// Block resource limits configuration structures.
//
// FRAME defines two resources that are limited within a block:
// - Weight (execution cost/time)
// - Length (block size)
//
// `frame_system` tracks consumption of each of these resources separately for each
// `DispatchClass`. This module contains configuration object for both resources,
// which should be passed to `frame_system` configuration when runtime is being set up.

// A ratio of `Normal` dispatch class within block, used as default value for
// `BlockWeight` and `BlockLength`. The `Default` impls are provided mostly for convenience
// to use in tests.

// ExtrinsicBaseWeight is the time to execute a NO-OP extrinsic, for example `System::remark`.
// Calculated by multiplying the *Average* with `1.0` and adding `0`.
//
// Stats nanoseconds:
//
//	Min, Max: 109_595, 114_170
//	Average:  110_536
//	Median:   110_233
//	Std-Dev:  933.39
//
// Percentiles nanoseconds:
//
//	99th: 114_120
//	95th: 112_680
//	75th: 110_858
var ExtrinsicBaseWeight types.Weight = types.WeightFromParts(WeightRefTimePerNanos.SaturatingMul(110_536), 0)

// Time to execute an empty block.
// Calculated by multiplying the *Average* with `1.0` and adding `0`.
//
// Stats nanoseconds:
//
//	Min, Max: 402_748, 458_228
//	Average:  412_772
//	Median:   406_151
//	Std-Dev:  13480.33
//
// Percentiles nanoseconds:
//
//	99th: 450_080
//	95th: 445_111
//	75th: 414_170
var BlockExecutionWeight types.Weight = types.WeightFromParts(WeightRefTimePerNanos.SaturatingMul(412_772), 0)

// MaximumBlockWeight is the maximum weight 2 seconds of compute with a 6 second average block time, with maximum proof size.
var MaximumBlockWeight types.Weight = types.WeightFromParts(WeightRefTimePerSecond.SaturatingMul(2), math.MaxUint64)

// DbWeight for RocksDB, used throughout the runtime.
var DbWeight types.RuntimeDbWeight = types.RuntimeDbWeight{
	Read:  25_000 * WeightRefTimePerNanos,
	Write: 100_000 * WeightRefTimePerNanos,
}

var WeightToFee types.WeightToFee = types.IdentityFee{}
var LengthToFee types.WeightToFee = types.IdentityFee{}
var FeeMultiplierUpdate types.WeightToFee = types.NewConstantMultiplier(sc.NewU128FromUint64(1))
