package types

import sc "github.com/LimeChain/goscale"

type IncRefStatus = sc.U8

const (
	IncRefStatusCreated IncRefStatus = iota
	IncRefStatusExisted
)
