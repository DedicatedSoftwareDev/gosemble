package support

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/hashing"
	"github.com/LimeChain/gosemble/primitives/storage"
)

type StorageValue[T sc.Encodable] struct {
	prefix     []byte
	name       []byte
	decodeFunc func(buffer *bytes.Buffer) T
}

func NewStorageValue[T sc.Encodable](prefix []byte, name []byte, decodeFunc func(buffer *bytes.Buffer) T) *StorageValue[T] {
	return &StorageValue[T]{
		prefix,
		name,
		decodeFunc,
	}
}

func (sv StorageValue[T]) Get() T {
	prefixHash := hashing.Twox128(sv.prefix)
	nameHash := hashing.Twox128(sv.name)

	return storage.GetDecode(append(prefixHash, nameHash...), sv.decodeFunc)
}

func (sv StorageValue[T]) Exists() bool {
	prefixHash := hashing.Twox128(sv.prefix)
	nameHash := hashing.Twox128(sv.name)

	exists := storage.Exists(append(prefixHash, nameHash...))

	return exists != 0
}

func (sv StorageValue[T]) Put(value T) {
	prefixHash := hashing.Twox128(sv.prefix)
	nameHash := hashing.Twox128(sv.name)

	storage.Set(append(prefixHash, nameHash...), value.Bytes())
}

func (sv StorageValue[T]) Take() []byte {
	prefixHash := hashing.Twox128(sv.prefix)
	nameHash := hashing.Twox128(sv.name)

	return storage.TakeBytes(append(prefixHash, nameHash...))
}

func (sv StorageValue[T]) TakeExact() T {
	prefixHash := hashing.Twox128(sv.prefix)
	nameHash := hashing.Twox128(sv.name)

	return storage.TakeDecode(append(prefixHash, nameHash...), sv.decodeFunc)
}
