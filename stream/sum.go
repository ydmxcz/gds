package stream

import (
	"github.com/ydmxcz/gds/util/constraints"
)

func Sum[T constraints.Integer](stm Stream[T]) (max T) {
	sumFunc := func(a, b T) T {
		return a + b
	}
	return FoldWith(stm, 0, sumFunc, sumFunc)
}
