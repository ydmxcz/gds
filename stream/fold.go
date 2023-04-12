package stream

import (
	"github.com/ydmxcz/gds/iterator"
)

func Fold[T any](stm Stream[T], init T, f func(T, T) T) (val T) {
	return FoldWith(stm, init, f, f)
}

func FoldWith[T, R any](stm Stream[T], init R, f func(R, T) R, combiner func(R, R) R) (val R) {
	if stm.parallelism == 0 {
		iterGenerators := stm.Activate(stm.parallelism)
		accum := init

		for {
			pull, b := iterGenerators()
			if !b {
				break
			}
			for val, ok := pull(); ok; val, ok = pull() {
				accum = f(accum, val)
			}
		}
		return accum
	} else {
		// TODO:parallel collect
		resChan := make(chan R, stm.parallelism)
		accum := init
		doParallel(stm.Activate(stm.parallelism),
			func(_ int, pull iterator.Iter[T]) {
				accum := init
				for val, ok := pull(); ok; val, ok = pull() {
					accum = f(accum, val)
				}
				resChan <- accum

				return
			})
		counter := 0
		for {
			if counter == stm.parallelism {
				break
			}
			accum = combiner(accum, <-resChan)
			counter++
		}
		close(resChan)
		return accum
	}
}
