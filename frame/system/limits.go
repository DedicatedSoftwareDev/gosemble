package system

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants"
	"github.com/LimeChain/gosemble/primitives/log"
	"github.com/LimeChain/gosemble/primitives/types"
)

type BlockLength struct {
	//  Maximal total length in bytes for each extrinsic class.
	//
	// In the worst case, the total block length is going to be:
	// `MAX(max)`
	Max types.PerDispatchClass[sc.U32]
}

func (bl BlockLength) Encode(buffer *bytes.Buffer) {
	bl.Max.Encode(buffer)
}

func (bl BlockLength) Bytes() []byte {
	return sc.EncodedBytes(bl)
}

func DefaultBlockLength() BlockLength {
	return MaxWithNormalRatio(
		constants.FiveMbPerBlockPerExtrinsic,
		constants.NormalDispatchRatio,
	)

	// return MaxWithNormalRatio(constants.FiveMbPerBlockPerExtrinsic, constants.DefaultNormalRatio)
}

// MaxWithNormalRatio Create new `BlockLength` with `max` for `Operational` & `Mandatory`
// and `normal * max` for `Normal`.
func MaxWithNormalRatio(max sc.U32, normal types.Perbill) BlockLength {
	return BlockLength{
		Max: types.PerDispatchClass[sc.U32]{
			Normal:      normal.Mul(max).(sc.U32),
			Operational: max,
			Mandatory:   max,
		},
	}
}

// WeightsPerClass `DispatchClass`-specific weight configuration.
type WeightsPerClass struct {
	// Base weight of single extrinsic of given class.
	BaseExtrinsic types.Weight

	// Maximal weight of single extrinsic. Should NOT include `base_extrinsic` cost.
	//
	// `None` indicates that this class of extrinsics doesn't have a limit.
	MaxExtrinsic sc.Option[types.Weight]

	// Block maximal total weight for all extrinsics of given class.
	//
	// `None` indicates that weight sum of this class of extrinsics is not
	// restricted. Use this value carefully, since it might produce heavily oversized
	// blocks.
	//
	// In the worst case, the total weight consumed by the class is going to be:
	// `MAX(max_total) + MAX(reserved)`.
	MaxTotal sc.Option[types.Weight]

	// Block reserved allowance for all extrinsics of a particular class.
	//
	// Setting to `None` indicates that extrinsics of that class are allowed
	// to go over total block weight (but at most `max_total` for that class).
	// Setting to `Some(x)` guarantees that at least `x` weight of particular class
	// is processed in every block.
	Reserved sc.Option[types.Weight]
}

func (cl WeightsPerClass) Encode(buffer *bytes.Buffer) {
	cl.BaseExtrinsic.Encode(buffer)
	cl.MaxExtrinsic.Encode(buffer)
	cl.MaxTotal.Encode(buffer)
	cl.Reserved.Encode(buffer)
}

func DecodeWeightsPerClass(buffer *bytes.Buffer) WeightsPerClass {
	cl := WeightsPerClass{}
	cl.BaseExtrinsic = types.DecodeWeight(buffer)
	cl.MaxExtrinsic = sc.DecodeOptionWith(buffer, types.DecodeWeight)
	cl.MaxTotal = sc.DecodeOptionWith(buffer, types.DecodeWeight)
	cl.Reserved = sc.DecodeOptionWith(buffer, types.DecodeWeight)
	return cl
}

func (cl WeightsPerClass) Bytes() []byte {
	return sc.EncodedBytes(cl)
}

type BlockWeights struct {
	// Base weight of block execution.
	BaseBlock types.Weight
	// Maximal total weight consumed by all kinds of extrinsics (without `reserved` space).
	MaxBlock types.Weight
	// Weight limits for extrinsics of given dispatch class.
	PerClass types.PerDispatchClass[WeightsPerClass]
}

func (bw BlockWeights) Encode(buffer *bytes.Buffer) {
	bw.BaseBlock.Encode(buffer)
	bw.MaxBlock.Encode(buffer)
	bw.PerClass.Encode(buffer)
}

func (bw BlockWeights) Bytes() []byte {
	return sc.EncodedBytes(bw)
}

func DefaultBlockWeights() BlockWeights {
	return WithSensibleDefaults(
		constants.MaximumBlockWeight,
		constants.NormalDispatchRatio,
	)
}

// Get per-class weight settings.
func (bw BlockWeights) Get(class types.DispatchClass) *WeightsPerClass {
	if class.Is(types.DispatchClassNormal) {
		return &bw.PerClass.Normal
	} else if class.Is(types.DispatchClassOperational) {
		return &bw.PerClass.Operational
	} else if class.Is(types.DispatchClassMandatory) {
		return &bw.PerClass.Mandatory
	} else {
		log.Critical("Invalid dispatch class")
	}

	panic("unreachable")
}

// WithSensibleDefaults Create a sensible default weights system given only expected maximal block weight and the
// ratio that `Normal` extrinsics should occupy.
//
// Assumptions:
//   - Average block initialization is assumed to be `10%`.
//   - `Operational` transactions have reserved allowance (`1.0 - normal_ratio`)
func WithSensibleDefaults(expectedBlockWeight types.Weight, normalRatio types.Perbill) BlockWeights {
	return NewBlockWeightsBuilder().
		BaseBlock(constants.BlockExecutionWeight).
		ForClass(types.DispatchClassAll(), func(weights *WeightsPerClass) {
			weights.BaseExtrinsic = constants.ExtrinsicBaseWeight
		}).
		ForClass([]types.DispatchClass{types.NewDispatchClassNormal()}, func(weights *WeightsPerClass) {
			weights.MaxTotal = sc.NewOption[types.Weight](constants.NormalDispatchRatio.Mul(constants.MaximumBlockWeight))
		}).
		ForClass([]types.DispatchClass{types.NewDispatchClassOperational()}, func(weights *WeightsPerClass) {
			weights.MaxTotal = sc.NewOption[types.Weight](constants.MaximumBlockWeight)
			// Operational transactions have some extra reserved space, so that they
			// are included even if block reached `MAXIMUM_BLOCK_WEIGHT`.
			weights.Reserved =
				sc.NewOption[types.Weight]((constants.MaximumBlockWeight.Sub(constants.NormalDispatchRatio.Mul(constants.MaximumBlockWeight).(types.Weight))))
		}).
		AvgBlockInitialization(constants.AverageOnInitializeRatio).
		Build()
	// TODO: builder.Expect("Sensible defaults are tested to be valid")
}

// An opinionated builder for `Weights` object.
type BlockWeightsBuilder struct {
	Weights  BlockWeights
	InitCost sc.Option[types.Perbill]
}

// Start constructing new `BlockWeights` object.
//
// By default all kinds except of `Mandatory` extrinsics are disallowed.
func NewBlockWeightsBuilder() *BlockWeightsBuilder {
	// Start constructing new `BlockWeights` object.
	//
	// By default all kinds except of `Mandatory` extrinsics are disallowed.
	WeightsForNormalAndOperational := WeightsPerClass{
		BaseExtrinsic: constants.ExtrinsicBaseWeight,
		MaxExtrinsic:  sc.NewOption[types.Weight](nil),
		MaxTotal:      sc.NewOption[types.Weight](types.WeightZero()),
		Reserved:      sc.NewOption[types.Weight](types.WeightZero()),
	}

	WeightsForMandatory := WeightsPerClass{
		BaseExtrinsic: constants.ExtrinsicBaseWeight,
		MaxExtrinsic:  sc.NewOption[types.Weight](nil),
		MaxTotal:      sc.NewOption[types.Weight](nil),
		Reserved:      sc.NewOption[types.Weight](nil),
	}

	weightsPerClass := types.PerDispatchClass[WeightsPerClass]{
		Mandatory:   WeightsForMandatory,
		Normal:      WeightsForNormalAndOperational,
		Operational: WeightsForNormalAndOperational,
	}

	return &BlockWeightsBuilder{
		Weights: BlockWeights{
			BaseBlock: constants.BlockExecutionWeight,
			MaxBlock:  types.WeightZero(),
			PerClass:  weightsPerClass,
		},
		InitCost: sc.NewOption[types.Perbill](nil),
	}
}

// Set base block weight.
func (b *BlockWeightsBuilder) BaseBlock(baseBlock types.Weight) *BlockWeightsBuilder {
	b.Weights.BaseBlock = baseBlock
	return b
}

// ForClass Set parameters for particular class.
//
// Note: `None` values of `max_extrinsic` will be overwritten in `build` in case
// `avg_block_initialization` rate is set to a non-zero value.
func (b *BlockWeightsBuilder) ForClass(classes []types.DispatchClass, action func(_ *WeightsPerClass)) *BlockWeightsBuilder {
	for _, cl := range classes {
		action(b.Weights.PerClass.Get(cl))
	}
	return b
}

// AvgBlockInitialization Average block initial ization weight cost.
//
// This value is used to derive maximal allowed extrinsic weight for each
// class, based on the allowance.
//
// This is to make sure that extrinsics don't stay forever in the pool,
// because they could seamingly fit the block (since they are below `max_block`),
// but the cost of calling `on_initialize` always prevents them from being included.
func (b *BlockWeightsBuilder) AvgBlockInitialization(initCost types.Perbill) *BlockWeightsBuilder {
	b.InitCost = sc.NewOption[types.Perbill](initCost)
	return b
}

// Construct the `BlockWeights` object.
func (b *BlockWeightsBuilder) Build() BlockWeights {
	// compute max extrinsic size
	weights, initCost := b.Weights, b.InitCost

	// compute max block size.
	for _, class := range types.DispatchClassAll() {
		if (*weights.PerClass.Get(class)).MaxTotal.HasValue {
			max := (*weights.PerClass.Get(class)).MaxTotal.Value
			weights.MaxBlock = max.Max(weights.MaxBlock)
		}
	}

	// compute max size of single extrinsic
	var initWeight sc.Option[types.Weight]
	if initCost.HasValue {
		initWeight = sc.NewOption[types.Weight](initCost.Value.Mul(weights.MaxBlock))
	} else {
		initWeight = sc.NewOption[types.Weight](nil)
	}

	if initWeight.HasValue {
		for _, class := range types.DispatchClassAll() {
			perClass := *(weights.PerClass.Get(class))
			if !perClass.MaxExtrinsic.HasValue && initCost.HasValue {
				if perClass.MaxTotal.HasValue {
					perClass.MaxExtrinsic = sc.NewOption[types.Weight](perClass.MaxTotal.Value.SaturatingSub(initWeight.Value).SaturatingSub(perClass.BaseExtrinsic))
				} else {
					perClass.MaxExtrinsic = sc.NewOption[types.Weight](nil)
				}
			}
		}
	}

	// Validate the result
	// TODO: weights.Validate()
	return weights
}
