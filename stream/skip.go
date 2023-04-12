package stream

import (
	"sync/atomic"

	"github.com/ydmxcz/gds/fn"
	"github.com/ydmxcz/gds/iterator"
)

func Skip[T any](stm Stream[T], s int) Stream[T] {
	generater := stm.Activate
	skip := int64(s)
	var curr int64 = 0
	return Stream[T]{
		parallelism: stm.parallelism,
		Activate: func(parallelism int) iterator.Iter[iterator.Iter[T]] {

			segementer := generater(parallelism)

			return func() (pr iterator.Iter[T], ok bool) {
				if pull, o1 := segementer(); o1 {
					return func() (val T, ok bool) {

						for val, ok := pull(); ok; val, ok = pull() {
							if atomic.AddInt64(&curr, 1) < skip {
								continue
							}
							return val, ok
						}
						return
					}, true
				}
				return nil, false

			}
		},
	}
}

func SkipWhile[T any](stm Stream[T], while fn.Predicate[T]) Stream[T] {
	generater := stm.Activate
	return Stream[T]{
		parallelism: stm.parallelism,
		Activate: func(parallelism int) iterator.Iter[iterator.Iter[T]] {

			segementer := generater(parallelism)

			return func() (pr iterator.Iter[T], ok bool) {
				if pull, o1 := segementer(); o1 {
					return func() (val T, ok bool) {

						for val, ok := pull(); ok; val, ok = pull() {
							if while(val) {
								continue
							}
							return val, ok
						}
						return
					}, true
				}
				return nil, false

			}
		},
	}
}
