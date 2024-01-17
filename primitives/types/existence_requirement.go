package types

import sc "github.com/LimeChain/goscale"

type ExistenceRequirement sc.U8

const (
	ExistenceRequirementKeepAlive ExistenceRequirement = iota
	ExistenceRequirementAllowDeath
)
