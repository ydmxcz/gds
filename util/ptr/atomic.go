package ptr

import (
	"sync/atomic"
	"unsafe"
)

func StorePointer[T any](p **T, new *T) {
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(p)), unsafe.Pointer(new))
}

func LoadPointer[T any](p **T) *T {
	return (*T)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(p))))
}

func SwapPointer[T any](p **T, newPtr *T) *T {
	return (*T)(atomic.SwapPointer((*unsafe.Pointer)(unsafe.Pointer(p)), unsafe.Pointer(newPtr)))
}

func CompareAndSwap[T any](p **T, old, new *T) (ok bool) {
	return atomic.CompareAndSwapPointer(
		(*unsafe.Pointer)(unsafe.Pointer(p)), unsafe.Pointer(old), unsafe.Pointer(new))
}
