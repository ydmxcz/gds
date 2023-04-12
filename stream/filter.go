package stream

import "github.com/ydmxcz/gds/iterator"

func Filter[T any](stm Stream[T], filt func(T) bool) Stream[T] {
	generater := stm.Activate

	return Stream[T]{
		Activate: func(parallelism int) iterator.Iter[iterator.Iter[T]] {

			segementer := generater(parallelism)

			return func() (pr iterator.Iter[T], ok bool) {
				if pull, ok := segementer(); ok {
					return func() (val T, ok bool) {
						for t, ok := pull(); ok; t, ok = pull() {
							if filt(t) {
								return t, true
							}
						}
						return val, false
					}, true
				}
				return nil, false

			}
		},
		parallelism: stm.parallelism,
	}
}
