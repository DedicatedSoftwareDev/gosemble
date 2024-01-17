package types

import sc "github.com/LimeChain/goscale"

type Module interface {
	Functions() map[sc.U8]Call
	PreDispatch(call Call) (sc.Empty, TransactionValidityError)
	ValidateUnsigned(source TransactionSource, call Call) (ValidTransaction, TransactionValidityError)
	Metadata() (sc.Sequence[MetadataType], MetadataModule)
}
