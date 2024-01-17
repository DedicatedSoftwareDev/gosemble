package types

import (
	sc "github.com/LimeChain/goscale"
)

const (
	// Commit the transaction.
	TransactionOutcomeCommit sc.U8 = iota
	// Rollback the transaction.
	TransactionOutcomeRollback
)

// TransactionOutcome Describes on what should happen with a storage transaction.
type TransactionOutcome = sc.VaryingData

func NewTransactionOutcomeCommit(res sc.Encodable) TransactionOutcome {
	return sc.NewVaryingData(TransactionOutcomeCommit, res)
}

func NewTransactionOutcomeRollback(res sc.Encodable) TransactionOutcome {
	return sc.NewVaryingData(TransactionOutcomeRollback, res)
}
