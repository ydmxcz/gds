package bplustree_test

import (
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/ydmxcz/gds/collections/bplustree"
	"github.com/ydmxcz/gds/fn"
)

var arr = []int{}

func initbptree() *bplustree.Tree[int, int] {
	rand.Seed(time.Now().UnixNano())
	tree := bplustree.NewMap[int, int](2, fn.Comp[int])
	for tree.Size() < 10000 {
		n := rand.Int()
		if tree.InsertUnique(n, n) {
			arr = append(arr, n)
		}
	}
	return tree
}

func initMap() map[int]int {
	rand.Seed(time.Now().UnixNano())
	tree := map[int]int{} //bplustree.NewMap[int, int](2, fn.Comp[int])
	for len(tree) < 10000 {
		n := rand.Int()
		if _, ok := tree[n]; !ok {
			tree[n] = n
			arr = append(arr, n)
		}
	}
	return tree
}

// func TestTree_Iterator(t *testing.T) {
// 	tree := initbptree()
// 	i := tree.Iter()
// 	count := 0
// 	for i.IsValid() {
// 		fmt.Println(i.GetValue())
// 		i.Next()
// 		count++
// 	}
// 	fmt.Println(count)
// 	i = tree.RIter()
// 	count = 0
// 	for i.IsValid() {
// 		fmt.Println(i.GetValue())
// 		i.Next()
// 		count++
// 	}
// 	fmt.Println(count)
// 	i = tree.End()
// 	i.Next()
// 	fmt.Println(i.GetValue(), i.IsValid())
// 	fmt.Println(tree.String())
// }

func BenchmarkTree_Put(b *testing.B) {
	tree := bplustree.NewMap[int, int](32, fn.Comp[int])
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tree.InsertUnique(i, i)
	}
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

// 10406180
// 2545668
func BenchmarkMapGet_Parallel(b *testing.B) {
	tree := initMap()
	lock := sync.Mutex{}
	b.ResetTimer()
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			for i := 0; i < len(arr); i++ {
				lock.Lock()
				_ = tree[arr[i]] //tree.Get(arr[i])
				lock.Unlock()
			}
		}

	})

}

func BenchmarkMap_Put(b *testing.B) {
	tree := map[int]int{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tree[i] = i
	}
}

func BenchmarkTree_Get(b *testing.B) {
	tree := bplustree.NewMap[int, int](32, fn.Comp[int])
	for i := 0; i < b.N; i++ {
		tree.InsertUnique(i, i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tree.Get(i)
	}
}

func BenchmarkMap_Get(b *testing.B) {
	tree := map[int]int{} //bplustree.NewMap[int, int](16, fn.Comp[int])
	for i := 0; i < b.N; i++ {
		tree[i] = i
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = tree[i]
	}
}
