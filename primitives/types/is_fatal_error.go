package types

import "github.com/LimeChain/goscale"

type IsFatalError interface {
	goscale.Encodable
	IsFatal() goscale.Bool
}
