package iterator

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"

	"github.com/ydmxcz/gds/fn"
)

type Iter[T any] fn.Pull[T]

type PushIter[T any] fn.Push[T]

type PushIterBool[T any] fn.PushPred[T]

type IterGenerator[T any] func(parallelism int) Iter[Iter[T]]

func YieldPrintln[T any](elem T) bool {
	fmt.Println(elem)
	return true
}

func CastToPush[T any](next Iter[T]) func(func(T) bool) {
	return func(yield func(T) bool) {
		for {
			v, ok := next()
			if !ok || !yield(v) {
				break
			}
		}
	}
}

func CastToPull[T any](push PushIter[T]) (next Iter[T], stop func()) {
	// var val V
	ctx, cf := context.WithCancel(context.Background())
	ch := make(chan T)
	go func() {
		push(func(v T) bool {
			select {
			case <-ctx.Done():
				return false
			case ch <- v:
				return true
			}
		})
		close(ch)
	}()
	return func() (val T, ok bool) {
		val, ok = <-ch
		return val, ok
	}, cf
}

func Lines(file string) func(func(string, error) bool) {
	return func(yield func(string, error) bool) {
		f, err := os.Open(file)
		if err != nil {
			yield("", err)
			return
		}
		defer f.Close()
		b := bufio.NewReader(f)
		for {
			line, err := b.ReadString('\n')
			if err != nil {
				if err != io.EOF {
					yield("", err)
				}
				break
			}
			if !yield(line, nil) {
				break
			}
		}
	}
}

func (next Iter[T]) All(yelid fn.Predicate[T]) {
	for v, ok := next(); ok; v, ok = next() {
		yelid(v)
	}
}

func FromSlice[T any](s []T) Iter[T] {
	var idx = 0
	return func() (val T, ok bool) {
		ok = idx < len(s)
		if ok {
			val = s[idx]
			idx++
		}
		return
	}
}

// func RangeSlice[T any](slice []T, yelid fn.Predicate[T]) {
// 	// for i := 0; i < len(slice); i++ {
// 	// 	if !yelid(slice[i]) {
// 	// 		return
// 	// 	}
// 	// }
// }

func FromChan[T any](c <-chan T) Iter[T] {
	return func() (val T, ok bool) {
		val, ok = <-c
		return val, ok
	}
}
