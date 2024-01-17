package types

import sc "github.com/LimeChain/goscale"

type Reasons sc.U8

const (
	ReasonsFee Reasons = iota
	ReasonsMisc
	ReasonsAll
)
