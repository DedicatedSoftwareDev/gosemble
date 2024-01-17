package dispatchables

import (
	"bytes"
	"math/big"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants"
	"github.com/LimeChain/gosemble/frame/balances/events"
	"github.com/LimeChain/gosemble/frame/system"
	"github.com/LimeChain/gosemble/primitives/hashing"
	"github.com/LimeChain/gosemble/primitives/storage"
	"github.com/LimeChain/gosemble/primitives/types"
)

type NegativeImbalance struct {
	types.Balance
}

func NewNegativeImbalance(balance types.Balance) NegativeImbalance {
	return NegativeImbalance{balance}
}

func (ni NegativeImbalance) Drop() {
	key := append(hashing.Twox128(constants.KeyBalances), hashing.Twox128(constants.KeyTotalIssuance)...)

	issuance := storage.GetDecode(key, sc.DecodeU128)

	issuanceBn := issuance.ToBigInt()
	sub := new(big.Int).Sub(issuanceBn, ni.ToBigInt())

	if sub.Cmp(issuanceBn) > 0 {
		sub = issuanceBn
	}

	storage.Set(key, sc.NewU128FromBigInt(sub).Bytes())
}

type PositiveImbalance struct {
	types.Balance
}

func NewPositiveImbalance(balance types.Balance) PositiveImbalance {
	return PositiveImbalance{balance}
}

func (pi PositiveImbalance) Drop() {
	key := append(hashing.Twox128(constants.KeyBalances), hashing.Twox128(constants.KeyTotalIssuance)...)

	issuance := storage.GetDecode(key, sc.DecodeU128)

	issuanceBn := issuance.ToBigInt()
	add := new(big.Int).Add(issuanceBn, pi.ToBigInt())

	if add.Cmp(issuanceBn) < 0 {
		add = issuanceBn
	}

	storage.Set(key, sc.NewU128FromBigInt(add).Bytes())
}

type DustCleanerValue struct {
	AccountId         types.Address32
	NegativeImbalance NegativeImbalance
}

func (dcv DustCleanerValue) Encode(buffer *bytes.Buffer) {
	dcv.AccountId.Encode(buffer)
	dcv.NegativeImbalance.Encode(buffer)
}

func (dcv DustCleanerValue) Bytes() []byte {
	return sc.EncodedBytes(dcv)
}

func (dcv DustCleanerValue) Drop() {
	system.DepositEvent(events.NewEventDustLost(dcv.AccountId.FixedSequence, dcv.NegativeImbalance.Balance))
	dcv.NegativeImbalance.Drop()
}
