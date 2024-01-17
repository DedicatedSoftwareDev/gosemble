//go:build nonwasmenv

package env

/*
	Miscellaneous: Interface that provides miscellaneous functions for communicating between the runtime and the node.
*/

func ExtMiscPrintHexVersion1(data int64) {
	panic("not implemented")
}

func ExtMiscPrintUtf8Version1(data int64) {
	panic("not implemented")
}

func ExtMiscRuntimeVersionVersion1(data int64) int64 {
	panic("not implemented")
}
