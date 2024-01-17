package events

import (
	"bytes"

	"github.com/LimeChain/gosemble/constants/balances"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/log"
	"github.com/LimeChain/gosemble/primitives/types"
)

// Balances module events.
const (
	EventEndowed sc.U8 = iota
	EventDustLost
	EventTransfer
	EventBalanceSet
	EventReserved
	EventUnreserved
	EventReserveRepatriated
	EventDeposit
	EventWithdraw
	EventSlashed
)

func NewEventEndowed(account types.PublicKey, freeBalance types.Balance) types.Event {
	return types.NewEvent(balances.ModuleIndex, EventEndowed, account, freeBalance)
}

func NewEventDustLost(account types.PublicKey, amount types.Balance) types.Event {
	return types.NewEvent(balances.ModuleIndex, EventDustLost, account, amount)
}

func NewEventTransfer(from types.PublicKey, to types.PublicKey, amount types.Balance) types.Event {
	return types.NewEvent(balances.ModuleIndex, EventTransfer, from, to, amount)
}

func NewEventBalanceSet(account types.PublicKey, free types.Balance, reserved types.Balance) types.Event {
	return types.NewEvent(balances.ModuleIndex, EventBalanceSet, account, free, reserved)
}

func NewEventReserved(account types.PublicKey, amount types.Balance) types.Event {
	return types.NewEvent(balances.ModuleIndex, EventReserved, account, amount)
}

func NewEventUnreserved(account types.PublicKey, amount types.Balance) types.Event {
	return types.NewEvent(balances.ModuleIndex, EventUnreserved, account, amount)
}

func NewEventReserveRepatriated(from types.PublicKey, to types.PublicKey, amount types.Balance, destinationStatus types.BalanceStatus) types.Event {
	return types.NewEvent(balances.ModuleIndex, EventReserveRepatriated, from, to, amount, destinationStatus)
}

func NewEventDeposit(account types.PublicKey, amount types.Balance) types.Event {
	return types.NewEvent(balances.ModuleIndex, EventDeposit, account, amount)
}

func NewEventWithdraw(account types.PublicKey, amount types.Balance) types.Event {
	return types.NewEvent(balances.ModuleIndex, EventWithdraw, account, amount)
}

func NewEventSlashed(account types.PublicKey, amount types.Balance) types.Event {
	return types.NewEvent(balances.ModuleIndex, EventSlashed, account, amount)
}

func DecodeEvent(buffer *bytes.Buffer) types.Event {
	module := sc.DecodeU8(buffer)
	if module != balances.ModuleIndex {
		log.Critical("invalid balances.Event module")
	}

	b := sc.DecodeU8(buffer)

	switch b {
	case EventEndowed:
		account := types.DecodePublicKey(buffer)
		freeBalance := sc.DecodeU128(buffer)
		return NewEventEndowed(account, freeBalance)
	case EventDustLost:
		account := types.DecodePublicKey(buffer)
		amount := sc.DecodeU128(buffer)
		return NewEventDustLost(account, amount)
	case EventTransfer:
		from := types.DecodePublicKey(buffer)
		to := types.DecodePublicKey(buffer)
		amount := sc.DecodeU128(buffer)
		return NewEventTransfer(from, to, amount)
	case EventBalanceSet:
		account := types.DecodePublicKey(buffer)
		free := sc.DecodeU128(buffer)
		reserved := sc.DecodeU128(buffer)
		return NewEventBalanceSet(account, free, reserved)
	case EventReserved:
		account := types.DecodePublicKey(buffer)
		amount := sc.DecodeU128(buffer)
		return NewEventReserved(account, amount)
	case EventUnreserved:
		account := types.DecodePublicKey(buffer)
		amount := sc.DecodeU128(buffer)
		return NewEventUnreserved(account, amount)
	case EventReserveRepatriated:
		from := types.DecodePublicKey(buffer)
		to := types.DecodePublicKey(buffer)
		amount := sc.DecodeU128(buffer)
		destinationStatus := types.DecodeBalanceStatus(buffer)
		return NewEventReserveRepatriated(from, to, amount, destinationStatus)
	case EventDeposit:
		account := types.DecodePublicKey(buffer)
		amount := sc.DecodeU128(buffer)
		return NewEventDeposit(account, amount)
	case EventWithdraw:
		account := types.DecodePublicKey(buffer)
		amount := sc.DecodeU128(buffer)
		return NewEventWithdraw(account, amount)
	case EventSlashed:
		account := types.DecodePublicKey(buffer)
		amount := sc.DecodeU128(buffer)
		return NewEventSlashed(account, amount)
	default:
		log.Critical("invalid balances.Event type")
	}

	panic("unreachable")
}
