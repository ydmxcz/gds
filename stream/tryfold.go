package stream

import (
	"context"

	"github.com/ydmxcz/gds/iterator"
)

func TryFold[T any](stm Stream[T], init T, f func(T, T) (T, bool)) (val T, ok bool) {
	if stm.parallelism == 0 {
		pull, b := stm.Activate(stm.parallelism)()
		if !b {
			return
		}
		accum := init
		for val, ok := pull(); ok; val, ok = pull() {
			accum, ok = f(accum, val)
			if !ok {
				return accum, false
			}
		}
		return accum, true
	} else {
		resChan := make(chan T, stm.parallelism)
		ctx, cancel := context.WithCancel(context.Background())
		accum := init
		doParallel(stm.Activate(stm.parallelism),
			func(_ int, pull iterator.Iter[T]) {
				accum := init
				for val, ok := pull(); ok; val, ok = pull() {
					select {
					case <-ctx.Done():
						return
					default:
						accum, ok = f(accum, val)
						if !ok {
							cancel()
							return
						}
					}
				}
				resChan <- accum
				return
			})
		counter := 0
		var ok bool
		for {
			if counter == stm.parallelism {
				break
			}
			select {
			case val = <-resChan:
				accum, ok = f(accum, val)
				if !ok {
					cancel()
					return accum, false
				}
				counter++
			case <-ctx.Done():
				return val, false
			}
		}
		return accum, true
	}
}
