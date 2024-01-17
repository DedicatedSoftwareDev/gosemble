//go:build !nonwasmenv

package env

/*
	Allocator: Provides functionality for calling into the memory allocator.
*/

//go:wasm-module env
//go:export ext_allocator_free_version_1
func ExtAllocatorFreeVersion1(ptr int32)

//go:wasm-module env
//go:export ext_allocator_malloc_version_1
func ExtAllocatorMallocVersion1(size int32) int32
