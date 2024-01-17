package types

import sc "github.com/LimeChain/goscale"

type DecRefStatus = sc.U8

const (
	DecRefStatusReaped sc.U8 = iota
	DecRefStatusExists
)
