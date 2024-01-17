//go:build !nonwasmenv

package env

/*
	Miscellaneous: Interface that provides miscellaneous functions for communicating between the runtime and the node.
*/

//go:wasm-module env
//go:export ext_misc_print_hex_version_1
func ExtMiscPrintHexVersion1(data int64)

//go:wasm-module env
//go:export ext_misc_print_utf8_version_1
func ExtMiscPrintUtf8Version1(data int64)

//go:wasm-module env
//go:export ext_misc_runtime_version_version_1
func ExtMiscRuntimeVersionVersion1(data int64) int64
