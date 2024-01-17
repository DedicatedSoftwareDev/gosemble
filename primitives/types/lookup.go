package types

import sc "github.com/LimeChain/goscale"

// AccountIdLookup A lookup implementation returning the `AccountId` from a `MultiAddress`.
type AccountIdLookup struct { // TODO: make it generic [AccountId, AccountIndex]
	// TODO: PhantomData[(AccountId, AccountIndex)]
}

func DefaultAccountIdLookup() AccountIdLookup {
	return AccountIdLookup{}
}

// TODO: MultiAddress[AccountId, AccountIndex]
func (l AccountIdLookup) Lookup(a MultiAddress) (Address32, TransactionValidityError) {
	address := LookupAddress(a)
	if address.HasValue {
		return address.Value, nil
	}

	return Address32{}, NewTransactionValidityError(NewUnknownTransactionCannotLookup())
}

// LookupAddress Lookup an address to get an Id, if there's one there.
func LookupAddress(a MultiAddress) sc.Option[Address32] { // TODO: MultiAddress[AccountId, AccountIndex]
	if a.IsAccountId() {
		return sc.NewOption[Address32](a.AsAccountId().Address32)
	}

	if a.IsAddress32() {
		return sc.NewOption[Address32](a.AsAddress32())
	}

	if a.IsAccountIndex() {
		return sc.NewOption[Address32](LookupIndex(a.AsAccountIndex()))
	}

	return sc.NewOption[Address32](nil)
}

// LookupIndex Lookup an T::AccountIndex to get an Id, if there's one there.
func LookupIndex(index AccountIndex) sc.Option[Address32] {
	// TODO:
	return sc.NewOption[Address32](nil)
}
