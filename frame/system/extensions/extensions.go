package system

import (
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/frame/transaction_payment"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

// TODO:
// we need to have a way for configuring any additional
// checks that need to be performed
//
// For example:
// CheckNonZeroSender
// CheckSpecVersion
// CheckTxVersion
// CheckGenesis
// CheckEra
// CheckNonce
// CheckWeight
// ChargeAssetTxPayment
//
// currently those checks are explicit, but
// depending on the configuration
// we could use reflection instead

type Extra primitives.SignedExtra

func (e Extra) AdditionalSigned() (ok primitives.AdditionalSigned, err primitives.TransactionValidityError) {
	ok = primitives.AdditionalSigned{} // FormatVersion: primitives.ExtrinsicFormatVersion

	specVersion, err := CheckSpecVersion{}.AdditionalSigned()
	if err != nil {
		return ok, err
	}
	ok.SpecVersion = specVersion

	transactionVersion, err := CheckTxVersion{}.AdditionalSigned()
	if err != nil {
		return ok, err
	}
	ok.TransactionVersion = transactionVersion

	genesisHash, err := CheckGenesis{}.AdditionalSigned()
	if err != nil {
		return ok, err
	}
	ok.GenesisHash = genesisHash

	blockHash, err := CheckMortality(e.Era).AdditionalSigned()
	if err != nil {
		return ok, err
	}
	ok.BlockHash = blockHash

	return ok, err
}

// Information on a transaction's validity and, if valid, on how it relates to other transactions.
func (e Extra) Validate(who *primitives.Address32, call *primitives.Call, info *primitives.DispatchInfo, length sc.Compact) (ok primitives.ValidTransaction, err primitives.TransactionValidityError) {
	valid := primitives.DefaultValidTransaction()

	ok, err = CheckNonZeroAddress(*who).Validate(who, call, info, length)
	if err != nil {
		return ok, err
	}
	valid = valid.CombineWith(ok)

	// TODO: CheckSpecVersion<Runtime>
	// TODO: CheckTxVersion<Runtime>
	// TODO: CheckGenesis<Runtime>

	ok, err = CheckMortality(e.Era).Validate(who, call, info, length)
	if err != nil {
		return ok, err
	}
	valid = valid.CombineWith(ok)

	ok, err = CheckNonce(e.Nonce).Validate(who, call, info, length)
	if err != nil {
		return ok, err
	}
	valid = valid.CombineWith(ok)

	ok, err = CheckWeight{}.Validate(who, call, info, length)
	if err != nil {
		return ok, err
	}
	valid = valid.CombineWith(ok)

	ok, err = transaction_payment.ChargeTransactionPayment(e.Fee).Validate(who, call, info, length)
	if err != nil {
		return ok, err
	}
	valid = valid.CombineWith(ok)

	return valid, err
}

func (e Extra) ValidateUnsigned(call *primitives.Call, info *primitives.DispatchInfo, length sc.Compact) (ok primitives.ValidTransaction, err primitives.TransactionValidityError) {
	valid := primitives.DefaultValidTransaction()

	ok, err = CheckWeight{}.ValidateUnsigned(call, info, length)
	if err != nil {
		return ok, err
	}
	valid = valid.CombineWith(ok)

	return valid, err
}

// Do any pre-flight stuff for a signed transaction.
//
// Make sure to perform the same checks as in [`Validate`].
func (e Extra) PreDispatch(who *primitives.Address32, call *primitives.Call, info *primitives.DispatchInfo, length sc.Compact) (ok primitives.Pre, err primitives.TransactionValidityError) {
	_, err = CheckNonZeroAddress(*who).PreDispatch(who, call, info, length)
	if err != nil {
		return ok, err
	}

	// TODO: CheckSpecVersion<Runtime>
	// TODO: CheckTxVersion<Runtime>
	// TODO: CheckGenesis<Runtime>

	_, err = CheckMortality(e.Era).PreDispatch(who, call, info, length)
	if err != nil {
		return ok, err
	}

	_, err = CheckNonce(e.Nonce).PreDispatch(who, call, info, length)
	if err != nil {
		return ok, err
	}

	_, err = CheckWeight{}.PreDispatch(who, call, info, length)
	if err != nil {
		return ok, err
	}

	pre, err := transaction_payment.ChargeTransactionPayment(e.Fee).PreDispatch(who, call, info, length)
	if err != nil {
		return ok, err
	}

	return pre, nil
}

func (e Extra) PreDispatchUnsigned(call *primitives.Call, info *primitives.DispatchInfo, length sc.Compact) (ok primitives.Pre, err primitives.TransactionValidityError) {
	_, err = CheckWeight{}.PreDispatchUnsigned(call, info, length)
	return ok, err
}

func (e Extra) PostDispatch(pre sc.Option[primitives.Pre], info *primitives.DispatchInfo, postInfo *primitives.PostDispatchInfo, length sc.Compact, result *primitives.DispatchResult) (primitives.Pre, primitives.TransactionValidityError) {
	_, err := CheckWeight{}.PostDispatch(pre, info, postInfo, length, result)

	_, err = transaction_payment.ChargeTransactionPayment(e.Fee).PostDispatch(pre, info, postInfo, length, result)
	if err != nil {
		return primitives.Pre{}, err
	}

	return primitives.Pre{}, err
}
