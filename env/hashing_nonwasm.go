//go:build nonwasmenv

package env

/*
	Hashing: Interface that provides functions for hashing with different algorithms.
*/

func ExtHashingBlake2128Version1(data int64) int32 {
	panic("not implemented")
}

func ExtHashingBlake2256Version1(data int64) int32 {
	panic("not implemented")
}

func ExtHashingKeccak256Version1(data int64) int32 {
	panic("not implemented")
}

func ExtHashingTwox128Version1(data int64) int32 {
	panic("not implemented")
}

func ExtHashingTwox64Version1(data int64) int32 {
	panic("not implemented")
}
