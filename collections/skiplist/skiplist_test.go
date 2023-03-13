package skiplist

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/ydmxcz/gds/fn"
)

//func TestSkipList_Insert(t *testing.T) {
//	arr := make([]int, 10)
//	rand.Seed(time.Now().UnixNano())
//	sl := New[int, string](comparison.Less[int])
//	for i := 0; i < 10; i++ {
//		arr[i] = i + 1 //rand.Intn(1000)
//	}
//	for i := 0; i < 10; i++ {
//		sl.InsertEquals(arr[i], strconv.Itoa(arr[i]))
//		// fmt.Println(insert)
//	}
//	sl.InsertUnique(arr[5], strconv.Itoa(arr[4]))
//	sl.InsertUnique(arr[5], strconv.Itoa(arr[4]))
//	sl.InsertUnique(arr[5], strconv.Itoa(arr[4]))
//
//	for i := 0; i < sl.Size(); i++ {
//		fmt.Print(sl.GetElementByRank(uint64(i+1)), ",")
//	}
//
//	fmt.Println()
//	fmt.Println(sl.Size(), ",OK")
//}
//
//func TestList_GetRank(t *testing.T) {
//	arr := make([]int, 10)
//	rand.Seed(time.Now().UnixNano())
//	sl := New[int, string](comparison.Less[int])
//	for i := 0; i < 10; i++ {
//		arr[i] = rand.Intn(1000)
//	}
//	for i := 0; i < 10; i++ {
//		sl.InsertEquals(arr[i], strconv.Itoa(arr[i]))
//	}
//	for i := 0; i < 10; i++ {
//		fmt.Print(sl.GetRank(arr[i]), ",")
//	}
//	fmt.Println()
//	fmt.Println(sl.Size(), ",OK")
//}
//
//func TestList_GetElementByRank(t *testing.T) {
//	arr := make([]int, 10)
//	rand.Seed(time.Now().UnixNano())
//	sl := New[int, string](comparison.Greater[int])
//	for i := 0; i < 10; i++ {
//		arr[i] = rand.Intn(1000)
//	}
//	for i := 0; i < 10; i++ {
//		sl.InsertEquals(arr[i], strconv.Itoa(arr[i]))
//	}
//	fmt.Println(arr[3])
//	sl.InsertEquals(arr[3], strconv.Itoa(arr[3]))
//	sl.InsertEquals(arr[3], strconv.Itoa(arr[3]))
//	for i := 1; i <= 10; i++ {
//		fmt.Print(sl.GetElementByRank(uint64(i)), ",")
//	}
//	fmt.Println()
//	fmt.Println()
//	fmt.Println()
//	fmt.Println()
//
//	fmt.Println(sl.String())
//	fmt.Println(sl.Size(), ",OK")
//}
//
//func arrtest(a *[10]int) {
//	a[6] = 666
//}
//
//func TestIsInRange(t *testing.T) {
//	// rand.Seed(time.Now().UnixNano())
//	sl := New[int, string](comparison.Less[int])
//	for i := 0; i < 10; i++ {
//		sl.InsertEquals(i, strconv.Itoa(i))
//	}
//	// The area in math is : [0,9],
//	// the util includes the area endpoint
//	fmt.Println(IsInAllRange(sl, 0, 9))   // true
//	fmt.Println(IsInAllRange(sl, -1, 9))  // false
//	fmt.Println(IsInAllRange(sl, -1, 10)) // false
//	fmt.Println(IsInAllRange(sl, -1, 0))  // false
//	fmt.Println(IsInAllRange(sl, 0, 10))  // false
//	fmt.Println(IsInAllRange(sl, 2, 9))   // true
//	fmt.Println(IsInAllRange(sl, 0, 1))   // true
//	fmt.Println(IsInAllRange(sl, 1, 1))   // false
//	fmt.Println("++++++++++++++++++++++++++++++++++")
//	fmt.Println(IsInPartOfRange(sl, 0, 9))   // true
//	fmt.Println(IsInPartOfRange(sl, -1, 9))  // true
//	fmt.Println(IsInPartOfRange(sl, -1, 10)) // true
//	fmt.Println(IsInPartOfRange(sl, -1, 0))  // true
//	fmt.Println(IsInPartOfRange(sl, 0, 10))  // true
//	fmt.Println(IsInPartOfRange(sl, 2, 9))   // true
//	fmt.Println(IsInPartOfRange(sl, 0, 1))   // true
//	fmt.Println(IsInPartOfRange(sl, 10, 13)) // false
//
//}
//
//func TestArrPtr(t *testing.T) {
//	var arr [10]int
//	for i := 0; i < 10; i++ {
//		arr[i] = i
//	}
//	arrtest(&arr)
//	fmt.Println(arr)
//}
//
//func TestComp(t *testing.T) {
//
//	fmt.Println(Comp(666, 667))
//	fmt.Println(Comp(667, 666))
//	fmt.Println(Comp(666, 666))
//	fmt.Println(Comp(666, 665))
//	fmt.Println(Comp(665, 666))
//}
//
//func Comp(a, b int) int {
//	if a > b {
//		return -1
//	}
//	if b > a {
//		return 1
//	}
//	return 0
//}
//
//func TestRange(t *testing.T) {
//	// rand.Seed(time.Now().UnixNano())
//	sl := New[int, string](comparison.Less[int])
//	for i := 0; i < 10; i++ {
//		sl.InsertEquals(i, strconv.Itoa(i))
//	}
//	fmt.Println(FirstInRange(sl, 0, 9).GetKey())
//	fmt.Println(LastInRange(sl, 0, 9).GetKey())
//
//	fmt.Println(FirstInRange(sl, 0, 10))
//	fmt.Println(LastInRange(sl, 0, 10))
//
//	fmt.Println(FirstInRange(sl, -1, 9))
//	fmt.Println(LastInRange(sl, -1, 9))
//
//	fmt.Println(FirstInRange(sl, -1, 10))
//	fmt.Println(LastInRange(sl, -1, 10))
//	for i := 1; i <= sl.Size(); i++ {
//		fmt.Print(sl.GetElementByRank(uint64(i)), " ")
//	}
//	fmt.Println()
//	fmt.Println(DeleteRangeByKey(sl, 3, 10))
//	fmt.Println()
//	for i := 1; i <= sl.Size(); i++ {
//		fmt.Print(sl.GetElementByRank(uint64(i)), " ")
//	}
//	fmt.Println()
//	fmt.Println(sl.Size())
//}

func constructTree(n int) *List[int, string] {
	rbTree := New[int, string](fn.Comp[int])
	for i := 0; i < n; i++ {
		rbTree.Put(i, strconv.Itoa(i))
	}
	return rbTree
}

func TestSkipListIter(t *testing.T) {
	skl := constructTree(10000000)
	iter := skl.Iter()
	count := 0
	for _, ok := iter(); ok; _, ok = iter() {
		// fmt.Print(v.Key, "  ")
		count++
	}
	fmt.Println(count)
}

func TestSkipListSplitableIter(t *testing.T) {
	l := constructTree(10)
	iter := l.SplitableIter()(2)
	for iter1, ok := iter(); ok; iter1, ok = iter() {
		for v1, o1 := iter1(); o1; v1, o1 = iter1() {
			fmt.Println(v1)
		}
		fmt.Println("====================")
	}
}

func BenchmarkList_Put(b *testing.B) {
	sl := New[int, int](fn.Comp[int])
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sl.Put(i, i)
	}
}

func BenchmarkList_Get(b *testing.B) {
	sl := New[int, int](fn.Comp[int])
	s := b.N
	for i := 0; i < s; i++ {
		sl.Put(i, i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sl.Get(i)
	}
}

type SafeSkipList[K comparable, V any] struct {
	tree  *List[K, V]
	mutex sync.RWMutex
}

func (st *SafeSkipList[K, V]) Put(key K, val V) bool {
	st.mutex.Lock()
	defer st.mutex.Unlock()
	return st.tree.Put(key, val)
}

func (st *SafeSkipList[K, V]) Get(key K) V {
	st.mutex.RLock()
	defer st.mutex.RUnlock()
	return st.tree.Get(key)
}

func BenchmarkSkipListRandomGet(b *testing.B) {
	num := 1024
	var sl = &SafeSkipList[int, int]{
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

type SafeMap[K comparable, V any] struct {
	tree  map[K]V //*redblacktree.TreeMap[K, V]
	mutex sync.RWMutex
}

func (st *SafeMap[K, V]) Put(key K, val V) bool {
	st.mutex.Lock()
	defer st.mutex.Unlock()
	// return st.tree.Put(key, val)
	st.tree[key] = val //st.tree.Put(key, val)
	return true
}

func (st *SafeMap[K, V]) Get(key K) V {
	st.mutex.RLock()
	defer st.mutex.RUnlock()
	return st.tree[key]
}

func BenchmarkMapGet(b *testing.B) {
	num := 1024
	var sl = &SafeMap[int, int]{
		tree: map[int]int{}, //redblacktree.NewTreeMap[int, int](comparison.Less[int]),
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

type SafeSyncMap[K comparable, V any] struct {
	tree  sync.Map //map[K]V //*redblacktree.TreeMap[K, V]
	mutex sync.RWMutex
}

func (st *SafeSyncMap[K, V]) Put(key K, val V) bool {
	// st.mutex.Lock()
	// defer st.mutex.Unlock()
	// return st.tree.Put(key, val)
	st.tree.Store(key, val) //st.tree.Put(key, val)
	return true
}

func (st *SafeSyncMap[K, V]) Get(key K) V {
	v, _ := st.tree.Load(key)
	return v.(V)
}

func BenchmarkSyncMapGet(b *testing.B) {
	num := 1024
	var sl = &SafeSyncMap[int, int]{
		// tree: map[int]int{}, //redblacktree.NewTreeMap[int, int](comparison.Less[int]),
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
