package utils

import (
	"unsafe"
)

var alivePointers = map[uintptr]interface{}{}

func Retain(data []byte) {
	ptr := &data[0]
	unsafePtr := uintptr(unsafe.Pointer(ptr))
	alivePointers[unsafePtr] = data
}

func Int64ToOffsetAndSize(offsetAndSize int64) (offset int32, size int32) {
	return int32(offsetAndSize), int32(offsetAndSize >> 32)
}

func OffsetAndSizeToInt64(offset int32, size int32) int64 {
	return int64(offset) | (int64(size) << 32)
}

func Offset32(data []byte) int32 {
	return int32(SliceToOffset(data))
}

func SliceToOffset(data []byte) uintptr {
	if len(data) == 0 {
		return uintptr(unsafe.Pointer(nil))
	}

	return uintptr(unsafe.Pointer(&data[0]))
}

func BytesToOffsetAndSize(data []byte) int64 {
	offset := SliceToOffset(data)
	size := len(data)
	return OffsetAndSizeToInt64(int32(offset), int32(size))
}

func ToWasmMemorySlice(offset int32, size int32) []byte {
	return unsafe.Slice((*byte)(unsafe.Pointer(uintptr(offset))), uintptr(size))
}
