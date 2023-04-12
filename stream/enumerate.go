package stream

import (
	"sync/atomic"

	"github.com/ydmxcz/gds/collections/truple"
	"github.com/ydmxcz/gds/iterator"
)

func Enumerate[T any](stm Stream[T]) Stream[truple.KV[int, T]] {
	generater := stm.Activate
	var curr int64 = 0
	return Stream[truple.KV[int, T]]{
		parallelism: stm.parallelism,
		Activate: func(parallelism int) iterator.Iter[iterator.Iter[truple.KV[int, T]]] {

			segementer := generater(parallelism)

			return func() (pr iterator.Iter[truple.KV[int, T]], ok bool) {
				if pull, o1 := segementer(); o1 {
					return func() (val truple.KV[int, T], ok bool) {
						if val, ok := pull(); ok {
							return truple.KV[int, T]{
								Key: int(atomic.AddInt64(&curr, 1)),
								Val: val,
							}, ok
						}

						return
					}, true
				}
				return nil, false

			}
		},
	}
}
