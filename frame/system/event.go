package system

import (
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/types"
)

// DepositEvent deposits an event into block's event record.
func DepositEvent(event types.Event) {
	depositEventIndexed([]types.H256{}, event)
}

// depositEventIndexed Deposits an event into this block's event record adding this event
// to the corresponding topic indexes.
//
// This will update storage entries that correspond to the specified topics.
// It is expected that light-clients could subscribe to this topics.
//
// NOTE: Events not registered at the genesis block and quietly omitted.
func depositEventIndexed(topics []types.H256, event types.Event) {
	blockNumber := StorageGetBlockNumber()
	if blockNumber == 0 {
		return
	}

	eventRecord := types.EventRecord{
		Phase:  StorageExecutionPhase(),
		Event:  event,
		Topics: topics,
	}

	oldEventCount := storageEventCount()
	newEventCount := oldEventCount + 1 // checked_add
	if newEventCount < oldEventCount {
		return
	}

	storageSetEventCount(newEventCount)

	storageAppendEvent(eventRecord)

	topicValue := sc.NewVaryingData(blockNumber, oldEventCount)
	for _, topic := range topics {
		storageAppendTopic(topic, topicValue)
	}
}
