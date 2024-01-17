package config

import (
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants/aura"
	"github.com/LimeChain/gosemble/constants/balances"
	"github.com/LimeChain/gosemble/constants/grandpa"
	"github.com/LimeChain/gosemble/constants/system"
	"github.com/LimeChain/gosemble/constants/testable"
	"github.com/LimeChain/gosemble/constants/timestamp"
	"github.com/LimeChain/gosemble/constants/transaction_payment"
	am "github.com/LimeChain/gosemble/frame/aura/module"
	bm "github.com/LimeChain/gosemble/frame/balances/module"
	gm "github.com/LimeChain/gosemble/frame/grandpa/module"
	sm "github.com/LimeChain/gosemble/frame/system/module"
	tm "github.com/LimeChain/gosemble/frame/testable/module"
	tsm "github.com/LimeChain/gosemble/frame/timestamp/module"
	tpm "github.com/LimeChain/gosemble/frame/transaction_payment/module"
	"github.com/LimeChain/gosemble/primitives/types"
)

// Modules contains all the modules used by the runtime.
var Modules = map[sc.U8]types.Module{
	system.ModuleIndex:              sm.NewSystemModule(),
	timestamp.ModuleIndex:           tsm.NewTimestampModule(),
	aura.ModuleIndex:                am.NewAuraModule(),
	grandpa.ModuleIndex:             gm.NewGrandpaModule(),
	balances.ModuleIndex:            bm.NewBalancesModule(),
	transaction_payment.ModuleIndex: tpm.NewTransactionPaymentModule(),
	testable.ModuleIndex:            tm.NewTestingModule(),
}
