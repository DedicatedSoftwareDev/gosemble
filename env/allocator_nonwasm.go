//go:build nonwasmenv

package env

/*
	Allocator: Provides functionality for calling into the memory allocator.
*/

func ExtAllocatorFreeVersion1(ptr int32) {
	panic("not implemented")
}

func ExtAllocatorMallocVersion1(size int32) int32 {
	panic("not implemented")
}
