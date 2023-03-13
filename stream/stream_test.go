package stream_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/ydmxcz/gds/collections/slice"
	"github.com/ydmxcz/gds/stream"
)

func TestFunc(t *testing.T) {
	sli := slice.Of(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1666, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1888, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1999, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	s := sli.Stream(8)
	stream.Collect(stream.Filter(stream.Map(s, func(a int) int {
		// time.Sleep(time.Millisecond * 1000)
		return a
	}), func(a int) bool {
		return a%2 == 0
	}), func(a int) bool {
		fmt.Println(a)
		return true
	})
}

func TestSum(t *testing.T) {
	sli := slice.Of(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	sum := stream.Fold(
		stream.Parallel(sli.Stream(), 8), 0, func(a, b int) int {
			// fmt.Println(a)
			time.Sleep(300 * time.Millisecond)
			return a + b
		})
	fmt.Println(sum)
}

func TestSI(t *testing.T) {
	sli := slice.Of(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1666, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1888, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1999, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	// ps := sli.SplitableIter()(4)

	ps := stream.Map(stream.New(sli.SplitableIter(), 4), func(a int) int {
		return a + 1000000
	}).Activate(4)
	// for pull, o1 := ps(); o1; pull, o1 = ps() {
	// 	for val, o2 := pull(); o2; val, o2 = pull() {
	// 		fmt.Println(val)
	// 	}
	// 	fmt.Println("=============================")
	// }
	k := 0
	for {
		pull, ok := ps()
		if !ok || k > 10 {
			fmt.Println("close chan", k)
			break
			// close(taskChan)
			// return
		}
		for val, o2 := pull(); o2; val, o2 = pull() {
			fmt.Println(val)
		}
		fmt.Println("=============================")
		k++
		// taskChan <- pull
		// fmt.Println("AAA")
	}
	fmt.Println("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	ps = stream.Map(stream.New(sli.SplitableIter(), 4), func(a int) int {
		return a + 1000000
	}).Activate(4)
	// for pull, o1 := ps(); o1; pull, o1 = ps() {
	// 	for val, o2 := pull(); o2; val, o2 = pull() {
	// 		fmt.Println(val)
	// 	}
	// 	fmt.Println("=============================")
	// }
	k = 0
	for {
		pull, ok := ps()
		if !ok || k > 10 {
			fmt.Println("close chan", k)
			break
			// close(taskChan)
			// return
		}
		for val, o2 := pull(); o2; val, o2 = pull() {
			fmt.Println(val)
		}
		fmt.Println("=============================")
		k++
		// taskChan <- pull
		// fmt.Println("AAA")
	}
}
