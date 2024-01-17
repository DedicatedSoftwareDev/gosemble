package balances

import sc "github.com/LimeChain/goscale"

const (
	ModuleIndex                    = sc.U8(4)
	FunctionTransferIndex          = 0
	FunctionSetBalanceIndex        = 1
	FunctionForceTransferIndex     = 2
	FunctionTransferKeepAliveIndex = 3
	FunctionTransferAllIndex       = 4
	FunctionForceFreeIndex         = 5
)
