package fn

import "github.com/ydmxcz/gds/util/constraints"

type Hashable interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64 | ~complex64 | ~complex128 | ~string
}

type Compare[T any] func(T, T) int

func Comp[T constraints.Ordered](v1, v2 T) int {
	if v1 > v2 {
		return 1
	} else if v1 < v2 {
		return -1
	} else {
		return 0
	}
}
