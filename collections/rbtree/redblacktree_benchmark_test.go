package rbtree

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/ydmxcz/gds/fn"
)

type SafeRbTree[K any, V any] struct {
	tree  *Tree[K, V]
	mutex sync.RWMutex
}

func (st *SafeRbTree[K, V]) Put(key K, val V) bool {
	st.mutex.Lock()
	defer st.mutex.Unlock()
	return st.tree.Put(key, val)
}

func (st *SafeRbTree[K, V]) Get(key K) V {
	st.mutex.RLock()
	defer st.mutex.RUnlock()
	return st.tree.Get(key)
}

func BenchmarkRandomGet(b *testing.B) {
	num := 102400
	var sl = &SafeRbTree[int, int]{
		tree: New[int, int](fn.Comp[int]),
	}
	rand.Seed(time.Now().UnixNano())
	nums := []int{}
	for i := 0; i < num; i++ {
		nums = append(nums, rand.Intn(num))
	}
	for i := 0; i < num; i++ {
		sl.Put(nums[i], 0)
	}
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			for i := 0; i < num; i++ {
				sl.Get(nums[i])
			}
		}
	})

}

func benchmarkRedBlackTreeInsert(b *testing.B) {
	rbTree := New[int, int](fn.Comp[int])
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rbTree.Put(i, 0)
	}
}

func benchmarkRedBlackTreeRemove(b *testing.B) {
	rbTree := New[int, int](fn.Comp[int])
	for i := 0; i < b.N; i++ {
		rbTree.Put(i, 0)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if rbTree.Delete(i) != nil {
			panic("Nil")
		}
	}

}

func benchmarkRedBlackTreeSearch(b *testing.B) {
	rbTree := New[int, int](fn.Comp[int])
	for i := 0; i < b.N; i++ {
		rbTree.Put(i, 0)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rbTree.Get(i)
	}
}

func BenchmarkRedBlackTree_Insert(b *testing.B) {
	rbTree := New[int, int](fn.Comp[int])
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rbTree.Put(i, 0)
	}
}

func BenchmarkRedBlackTree_Remove(b *testing.B) {
	for i := 0; i < 1; i++ {
		b.Run(strconv.Itoa(i), benchmarkRedBlackTreeRemove)

	}
}

func BenchmarkRedBlackTree_Search(b *testing.B) {
	for i := 0; i < 1; i++ {
		b.Run(strconv.Itoa(i), benchmarkRedBlackTreeSearch)
	}
	// b.Run("RemoveEquals", benchmarkRedBlackTreeRemove)
	// b.Run("search", benchmarkRedBlackTreeSearch)
}

var arr = []int{}

func initbptree() *Tree[int, int] {
	rand.Seed(time.Now().UnixNano())
	tree := New[int, int](fn.Comp[int])
	for tree.Size() < 10000 {
		n := rand.Int()
		if tree.Put(n, n) {
			arr = append(arr, n)
		}
	}
	return tree
}

func BenchmarkBPTreeGet_Parallel(b *testing.B) {
	tree := initbptree()
	lock := sync.Mutex{}
	b.ResetTimer()
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			for i := 0; i < len(arr); i++ {
				lock.Lock()
				tree.Get(arr[i])
				lock.Unlock()
			}
		}
	})
}
