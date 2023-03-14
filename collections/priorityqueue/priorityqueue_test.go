package priorityqueue_test

import (
	"fmt"
	"testing"

	"github.com/ydmxcz/gds/collections"
	"github.com/ydmxcz/gds/collections/priorityqueue"
	"github.com/ydmxcz/gds/fn"
)

func TestQueue_Push(t *testing.T) {
	q := priorityqueue.New(fn.Comp[int])
	arr := []int{653, 214, 386, 402, 518, 688, 778, 888, 999}
	for i := 0; i < len(arr); i++ {
		q.Push(arr[i])
	}
	fmt.Println(q.Peek())
	fmt.Println(q)
	fmt.Println(q.Peek())
	for i := 0; i < len(arr); i++ {
		//fmt.Println(q.Pop())
		q.Pop()
		fmt.Println(q.Peek())
	}
	//for i := 50; i > 0; i-- {
	//	q.Push(i)
	//}
}

func TestBuildHeapByOrderElement(t *testing.T) {
	q := priorityqueue.OfWith(fn.Comp[int], 653, 214, 386, 402, 518, 688, 778, 888, 999)
	fmt.Println(q)
	var qq collections.Queue[int] = q
	fmt.Println(qq.Len())
	//q1 := BuildHeapByOrderElement[int](util.Less[int],653,214,386,402,518,688,778,888,999)
	//fmt.Println(q1)
	//var qq2 sequence.Queue[int]  = q1
	//fmt.Println(qq2.Size())
}
