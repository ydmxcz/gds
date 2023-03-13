package stream

import (
	"sync"

	"github.com/ydmxcz/gds/fn"
	"github.com/ydmxcz/gds/iterator"
)

type IterGenerater[T any] func(parallelism int) iterator.Iter[iterator.Iter[T]]

type Stream[T any] struct {
	Activate    IterGenerater[T]
	parallelism int
}

func New[T any](generater IterGenerater[T], parallelism int) Stream[T] {
	return Stream[T]{Activate: generater, parallelism: parallelism}
}

func Map[T, R any](stm Stream[T], mapF func(T) R) Stream[R] {
	generater := stm.Activate

	return Stream[R]{
		parallelism: stm.parallelism,
		Activate: func(parallelism int) iterator.Iter[iterator.Iter[R]] {

			segementer := generater(parallelism)

			return func() (pr iterator.Iter[R], ok bool) {
				if pull, o1 := segementer(); o1 {
					return func() (val R, ok bool) {
						if t, ok := pull(); ok {
							return mapF(t), ok
						}
						return
					}, true
				}
				return nil, false

			}
		},
	}
}

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

func doParallel[T any](iterGenerater iterator.Iter[iterator.Iter[T]],
	parallelism int, yield fn.BinPredicate[int, iterator.Iter[T]]) *sync.WaitGroup {
	// TODO:parallel collect
	taskChan := make(chan iterator.Iter[T], parallelism)
	// res := make([][]T, stm.parallelism)
	wg := &sync.WaitGroup{}
	wg.Add(parallelism)

	go func(iGen iterator.Iter[iterator.Iter[T]], tChan chan iterator.Iter[T]) {
		for {
			pull, ok := iGen()
			if !ok {
				close(tChan)
				return
			}
			tChan <- pull
		}
	}(iterGenerater, taskChan)

	for i := 0; i < parallelism; i++ {

		go func(idx int, wg *sync.WaitGroup,
			tChan chan iterator.Iter[T]) {
			// r := make([]T, 0)
			pull, ok := <-tChan
			if !ok {
				wg.Done()
				return
			}
			yield(idx, pull)
			wg.Done()
		}(i, wg, taskChan)
	}
	return wg
}

func Collect[T any](stm Stream[T], collectTo fn.Predicate[T]) {
	if stm.parallelism == 0 {
		pull, b := stm.Activate(stm.parallelism)()
		if !b {
			return
		}
		for val, ok := pull(); ok; val, ok = pull() {
			collectTo(val)
		}
	} else {
		// TODO:parallel collect
		res := make([][]T, stm.parallelism)
		doParallel(stm.Activate(stm.parallelism), stm.parallelism,
			func(groupID int, pull iterator.Iter[T]) bool {
				r := make([]T, 0)
				for val, ok := pull(); ok; val, ok = pull() {
					r = append(r, val)
				}
				res[groupID] = r
				return true
			}).Wait()
		for i := 0; i < len(res); i++ {
			for j := 0; j < len(res[i]); j++ {
				collectTo(res[i][j])
			}
		}
	}
}

func Fold[T any](stm Stream[T], init T, f func(T, T) T) (val T) {
	if stm.parallelism == 0 {
		pull, b := stm.Activate(stm.parallelism)()
		if !b {
			return
		}
		accum := init
		for val, ok := pull(); ok; val, ok = pull() {
			accum = f(accum, val)
		}
		return accum
	} else {
		// TODO:parallel collect
		// resChan := make([][]T, stm.parallelism)
		resChan := make(chan T, stm.parallelism)
		accum := init
		doParallel(stm.Activate(stm.parallelism), stm.parallelism,
			func(_ int, pull iterator.Iter[T]) bool {
				accum := init
				for val, ok := pull(); ok; val, ok = pull() {
					accum = f(accum, val)
				}
				resChan <- accum
				// res[groupID] = r
				return true
			})
		counter := 0
		for {
			if counter == stm.parallelism {
				break
			}
			accum = f(accum, <-resChan)
			counter++
		}
		close(resChan)
		return accum
	}
}

func First[T ~int](stm Stream[T]) (T, bool) {
	var init T
	pull, b := stm.Activate(stm.parallelism)()
	if !b {
		return init, false
	}
	// accum := init
	for val, ok := pull(); ok; val, ok = pull() {
		return val, true
	}
	return init, false
}

func Max[T ~int](stm Stream[T]) (max T) {
	max, ok := First(stm)
	if !ok {
		return
	}
	return TryFold(stm, max, func(a, b T) (T, bool) {
		if a > b {
			return a, true
		}
		return b, false
	})
}

func Min[T ~int](stm Stream[T]) (min T) {
	min, ok := First(stm)
	if !ok {
		return
	}
	return TryFold(stm, min, func(a, b T) (T, bool) {
		if a < b {
			return a, true
		}
		return b, false
	})
}

func TryFold[T any](stm Stream[T], init T, f func(T, T) (T, bool)) (val T) {
	if stm.parallelism == 0 {
		pull, b := stm.Activate(stm.parallelism)()
		if !b {
			return
		}
		accum := init
		for val, ok := pull(); ok; val, ok = pull() {
			accum, ok = f(accum, val)
			if !ok {
				return accum
			}
		}
		return accum
	} else {
		// TODO:parallel collect
		// resChan := make([][]T, stm.parallelism)
		resChan := make(chan T, stm.parallelism)
		accum := init
		doParallel(stm.Activate(stm.parallelism), stm.parallelism,
			func(_ int, pull iterator.Iter[T]) bool {
				accum := init
				for val, ok := pull(); ok; val, ok = pull() {
					accum, ok = f(accum, val)
					if !ok {
						resChan <- accum
						return true
					}
				}
				resChan <- accum
				// res[groupID] = r
				return true
			})
		counter := 0
		var ok bool
		for {
			if counter == stm.parallelism {
				break
			}
			accum, ok = f(accum, val)
			if !ok {
				break
			}
			counter++
		}
		close(resChan)
		return accum
	}
}

func Parallel[T any](stm Stream[T], parallelism int) Stream[T] {
	stm.parallelism = parallelism
	return stm
}
