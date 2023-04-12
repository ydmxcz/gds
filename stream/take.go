package stream

import (
	"sync/atomic"

	"github.com/ydmxcz/gds/iterator"
)

func Take[T any](stm Stream[T], length int) Stream[T] {
	generater := stm.Activate
	l := int64(length)
	var curr int64 = 0
	return Stream[T]{
		parallelism: stm.parallelism,
		Activate: func(parallelism int) iterator.Iter[iterator.Iter[T]] {

			segementer := generater(parallelism)

			return func() (pr iterator.Iter[T], ok bool) {
				if pull, o1 := segementer(); o1 {
					return func() (val T, ok bool) {
						if atomic.LoadInt64(&curr) < l {
							if val, ok := pull(); ok {
								atomic.AddInt64(&curr, 1)
								return val, ok
							}
						}

						return
					}, true
				}
				return nil, false

			}
		},
	}
}
