package stream

import "github.com/ydmxcz/gds/iterator"

func FlatMap[T, U any](stm Stream[T], flatMap func(a T) Stream[U]) Stream[U] {
	generater := stm.Activate

	return Stream[U]{
		parallelism: stm.parallelism,
		Activate: func(parallelism int) iterator.Iter[iterator.Iter[U]] {

			segementer := generater(parallelism)

			return func() (pr iterator.Iter[U], ok bool) {
				if pull, o1 := segementer(); o1 {
					var iter iterator.Iter[U]
					var ok2 bool
					// call current iterator
					if t, ok := pull(); ok {
						// if successful,flat it.
						s := flatMap(t)
						iter, ok2 = s.Activate(0)()
						if !ok2 {
							// flat failed,make `iter` is nil
							iter = nil
						}
					}
					return func() (val U, ok bool) {
						// if iter is nil,no useable element
						if iter == nil {
							return
						}
						for {
							// get element from the current iterator
							if val, ok = iter(); ok {
								return val, true
							} else {
								// if failed,try to update current iterator
								if t, ok := pull(); ok {
									// everything sub-stream is signal thread,
									// parallel process just in the level of the current stream.
									iter, ok2 = flatMap(t).Activate(0)()
									if !ok2 {
										// update fail,return false
										return val, false
									}
									// update success,follow the for-loop to re-get a new element.
								} else {
									// no iterator anymore.
									return val, false
								}
							}
						}
					}, true
				}
				return nil, false

			}
		},
	}
}
