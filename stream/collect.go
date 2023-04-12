package stream

import (
	"github.com/ydmxcz/gds/fn"
	"github.com/ydmxcz/gds/iterator"
)

func Collect[T any](stm Stream[T], collectTo fn.Consumer[T]) {
	if stm.parallelism == 0 {
		iterGenerators := stm.Activate(stm.parallelism)
		for {
			pull, b := iterGenerators()
			if !b {
				return
			}
			for val, ok := pull(); ok; val, ok = pull() {
				collectTo(val)
			}
		}
	} else {

		doParallel(stm.Activate(stm.parallelism),
			func(groupID int, pull iterator.Iter[T]) {

				for val, ok := pull(); ok; val, ok = pull() {
					collectTo(val)
				}

			}).Wait()
	}
}
