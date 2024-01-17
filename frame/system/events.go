package system

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants/system"
	"github.com/LimeChain/gosemble/primitives/log"
	"github.com/LimeChain/gosemble/primitives/types"
)

// System module events.
const (
	EventExtrinsicSuccess sc.U8 = iota
	EventExtrinsicFailed
	EventCodeUpdated
	EventNewAccount
	EventKilledAccount
	EventRemarked
)

func NewEventExtrinsicSuccess(dispatchInfo types.DispatchInfo) types.Event {
	return types.NewEvent(system.ModuleIndex, EventExtrinsicSuccess, dispatchInfo)
}

func NewEventExtrinsicFailed(dispatchError types.DispatchError, dispatchInfo types.DispatchInfo) types.Event {
	return types.NewEvent(system.ModuleIndex, EventExtrinsicFailed, dispatchError, dispatchInfo)
}

func NewEventCodeUpdated() types.Event {
	return types.NewEvent(system.ModuleIndex, EventCodeUpdated)
}

func NewEventNewAccount(account types.PublicKey) types.Event {
	return types.NewEvent(system.ModuleIndex, EventNewAccount, account)
}

func NewEventKilledAccount(account types.PublicKey) types.Event {
	return types.NewEvent(system.ModuleIndex, EventKilledAccount, account)
}

func NewEventRemarked(sender types.PublicKey, hash types.H256) types.Event {
	return types.NewEvent(system.ModuleIndex, EventRemarked, sender, hash)
}

func DecodeEvent(buffer *bytes.Buffer) types.Event {
	moduleIndex := sc.DecodeU8(buffer)
	if moduleIndex != system.ModuleIndex {
		log.Critical("invalid system.Event")
	}

	b := sc.DecodeU8(buffer)

	switch b {
	case EventExtrinsicSuccess:
		dispatchInfo := types.DecodeDispatchInfo(buffer)
		return NewEventExtrinsicSuccess(dispatchInfo)
	case EventExtrinsicFailed:
		dispatchErr := types.DecodeDispatchError(buffer)
		dispatchInfo := types.DecodeDispatchInfo(buffer)
		return NewEventExtrinsicFailed(dispatchErr, dispatchInfo)
	case EventCodeUpdated:
		return NewEventCodeUpdated()
	case EventNewAccount:
		account := types.DecodePublicKey(buffer)
		return NewEventNewAccount(account)
	case EventKilledAccount:
		account := types.DecodePublicKey(buffer)
		return NewEventKilledAccount(account)
	case EventRemarked:
		account := types.DecodePublicKey(buffer)
		hash := types.DecodeH256(buffer)
		return NewEventRemarked(account, hash)
	default:
		log.Critical("invalid system.Event type")
	}

	panic("unreachable")
}
