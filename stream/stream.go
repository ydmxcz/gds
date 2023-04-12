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

func doParallel[T any](iterGenerater iterator.Iter[iterator.Iter[T]],
	yield fn.BinConsumer[int, iterator.Iter[T]]) *sync.WaitGroup {
	wg := &sync.WaitGroup{}
	idx := 0
	for {
		iter, ok := iterGenerater()
		if !ok {
			break
		}
		wg.Add(1)
		go func(idx int, wg *sync.WaitGroup, iter iterator.Iter[T]) {
			yield(idx, iter)
			wg.Done()
		}(idx, wg, iter)
		idx++
	}
	return wg
}

func Parallel[T any](stm Stream[T], parallelism int) Stream[T] {
	stm.parallelism = parallelism
	return stm
}
