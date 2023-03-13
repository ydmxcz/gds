package queue

import (
	"testing"
)

func getNodes(count int) []*Node[int] {
	cases := make([]*Node[int], 0, count)
	for i := 0; i < count; i++ {
		cases = append(cases, NewNode[int](i+100))
	}
	return cases
}

// func getNodeLs(count int) []*NodeL[int] {
// 	cases := make([]*NodeL[int], 0, count)
// 	for i := 0; i < count; i++ {
// 		cases = append(cases, NewNodeL[int](i+100))
// 	}
// 	return cases
// }

// func TestQueueBaseL(t *testing.T) {
// 	count := 5
// 	cases := getNodes(count)
// 	q := NewLockFree[int]()
// 	q.Push(cases[0])
// 	q.Push(cases[1])
// 	q.Push(cases[2])
// 	q.Push(cases[3])
// 	q.Push(cases[4])

// 	// for i := 0; i < count; i++ {
// 	// 	q.Push(cases[i])
// 	// 	if q.Len() != int64(i+1) {
// 	// 		t.Fatal("queue length error while pushing")
// 	// 	}
// 	// }
// 	fmt.Println(q.Len())
// 	fmt.Println(q.Pop())
// 	fmt.Println(q.Pop())
// 	fmt.Println(q.Pop())
// 	fmt.Println(q.Pop())
// 	fmt.Println(q.Pop())
// 	fmt.Println(q.Len())
// 	fmt.Println(q.Pop())
// 	// for i := 0; i < count; i++ {
// 	// 	q.Pop()
// 	// 	if q.Len() != int64(count-i-1) {
// 	// 		// fmt.Println(q.Len(), count-i-1)
// 	// 		t.Fatal("queue length error while poping")
// 	// 	}
// 	// }
// 	// if q.Pop() != nil {
// 	// 	t.Fatal("not nil")
// 	// }
// }

// func TestQueueBase(t *testing.T) {
// 	count := 10000
// 	cases := getNodes(count)
// 	q := NewLockFree[int]()
// 	for i := 0; i < count; i++ {
// 		q.Push(cases[i])
// 		if q.Len() != int64(i+1) {
// 			t.Fatal("queue length error while pushing")
// 		}
// 	}
// 	q.Len()
// 	for i := 0; i < count; i++ {
// 		q.Pop()
// 		if q.Len() != int64(count-i-1) {
// 			// fmt.Println(q.Len(), count-i-1)
// 			t.Fatal("queue length error while poping")
// 		}
// 	}
// 	if q.Pop() != nil {
// 		t.Fatal("not nil")
// 	}
// }

// func TestQueueCricle(t *testing.T) {
// 	count := 10000
// 	cases := getNodes(count)
// 	fmt.Println("AAAA")
// 	q := NewLockFree[int]()
// 	for i := 0; i < count; i++ {
// 		q.Push(cases[i])
// 	}
// 	fmt.Println("AAAA")
// 	q2 := NewLockFree[int]()
// 	for i := 0; i < count; i++ {
// 		// fmt.Println("B1")
// 		// n :=
// 		q2.Push(q.Pop())
// 		// fmt.Println("B2")

// 	}
// 	fmt.Println("AAAA")

// 	for i := 0; i < count; i++ {
// 		fmt.Println(q2.Pop().val)
// 	}
// }

func BenchmarkQueueCricle(b *testing.B) {
	count := 10000
	cases := getNodes(count)
	q := &Linked[int]{}
	for i := 0; i < count; i++ {
		q.PushNode(cases[i])
	}
	// q2 := NewQueue[int]()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		q.PushNode(q.PopNode())
	}
}

func BenchmarkQueueCricle_Parallel(b *testing.B) {
	count := 10000
	cases := getNodes(count)
	q := &Linked[int]{}
	for i := 0; i < count; i++ {
		q.PushNode(cases[i])
	}
	// fmt.Println(q.Len())
	b.ResetTimer()
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			q.PushNode(q.PopNode())

		}
	})
	// for i := 0; i < b.N; i++ {
	// 	q.Push(q.Pop())
	// }
}

// func TestQueueParllel(t *testing.T) {
// 	lq := NewLockFree[int]()
// 	wg := sync.WaitGroup{}
// 	wg.Add(10)
// 	for i := 0; i < 10; i++ {
// 		go func(id int, lq *LockFree[int], w *sync.WaitGroup) {
// 			for j := 0; j < 100; j++ {
// 				lq.Push(NewNodeL((id * 100) + j))
// 			}
// 			w.Done()
// 		}(i+1, lq, &wg)
// 	}
// 	wg.Wait()
// 	fmt.Println("Len:", lq.Len())
// 	wg.Add(10)
// 	for i := 0; i < 10; i++ {
// 		go func(id int, lq *LockFree[int], w *sync.WaitGroup) {
// 			for !lq.Empty() {
// 				fmt.Println(id, "::", lq.Pop()) //Push((id * 100) + j)
// 			}
// 			w.Done()
// 		}(i+1, lq, &wg)
// 	}
// 	wg.Wait()
// 	fmt.Println(lq.Len())
// }
