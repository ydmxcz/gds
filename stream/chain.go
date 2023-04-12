package stream

import (
	"fmt"

	"github.com/ydmxcz/gds/iterator"
)

func Chain[T any](stm ...Stream[T]) Stream[T] {

	return Stream[T]{
		Activate: func(parallelism int) iterator.Iter[iterator.Iter[T]] {
			segementers := make([]iterator.Iter[iterator.Iter[T]], 0, len(stm))
			for i := 0; i < len(stm); i++ {
				segementers = append(segementers, stm[i].Activate(stm[i].parallelism))
			}
			return func() (pr iterator.Iter[T], ok bool) {
				pulls := make([]iterator.Iter[T], 0, len(segementers))
				for i := 0; i < len(segementers); i++ {
					if pull, ok := segementers[i](); ok {
						pulls = append(pulls, pull)
					}
				}
				if len(pulls) != 0 {
					idx := 0
					pull := pulls[idx]
					return func() (val T, ok bool) {
						for {
							val, ok = pull()
							if ok {
								return val, true
							} else {
								idx++
								if idx == len(pulls) {
									fmt.Println("KKK")
									break
								}
								pull = pulls[idx]
							}
						}
						return val, false
					}, true

				}
				// if pull, ok := segementer(); ok {
				// 	pull2, _ := segementer2()
				// }

				return nil, false

			}
		},
		parallelism: 0,
	}
}
