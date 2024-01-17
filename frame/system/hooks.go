package system

import (
	"github.com/LimeChain/gosemble/primitives/types"
)

func onCreatedAccount(who types.Address32) {
	// hook on creating new account, currently not used in Substrate
	//T::OnNewAccount::on_new_account(&who);
	DepositEvent(NewEventNewAccount(who.FixedSequence))
}

func onKilledAccount(who types.Address32) {
	DepositEvent(NewEventKilledAccount(who.FixedSequence))
}
