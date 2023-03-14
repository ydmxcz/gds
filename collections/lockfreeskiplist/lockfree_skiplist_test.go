package lockfreeskiplist

import (
	"flag"
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/ydmxcz/gds/fn"
)

var n int

func init() {
	flag.IntVar(&n, "n", 10000, "element count")
}

func assert(b *testing.B, cond bool, message string) {
	if cond {
		return
	}
	b.Fatal(message)
}

func TestSkipList(t *testing.T) {
	var sl = NewLockFreeSkipList(fn.Comp[int])
	sl.Add(16)
	sl.Add(32)
	sl.Add(64)
	sl.Add(128)
	fmt.Println(sl.Get(15), sl.Get(17), sl.Get(16))
	fmt.Println(sl.Get(30), sl.Get(40), sl.Get(32))
	fmt.Println(sl.Get(60), sl.Get(70), sl.Get(64))
	fmt.Println(sl.Get(120), sl.Get(130), sl.Get(128))
}

func TestSkipListRank(t *testing.T) {
	var sl = NewLockFreeSkipList(fn.Comp[int])
	now := time.Now().UnixMilli()
	fmt.Println(now, (now - now%(3)))
	sl.Add(16)
	sl.Add(32)
	sl.Add(64)
	sl.Add(128)
	fmt.Println("size:", sl.Len())
	fmt.Println("rank:", sl.Rank(128))
	// fmt.Println(sl.GetRank(15), sl.GetRank(17), sl.GetRank(16))
	// fmt.Println(sl.GetRank(30), sl.GetRank(40), sl.GetRank(32))
	// fmt.Println(sl.GetRank(60), sl.GetRank(70), sl.GetRank(64))
	// fmt.Println(sl.GetRank(120), sl.GetRank(130), sl.GetRank(128))
}

// Add n elements per goroutines.
func BenchmarkRandomAdd(b *testing.B) {
	var sl = NewLockFreeSkipList(fn.Comp[int])

	var count int32
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if sl.Add(rand.Int() % n) {
				atomic.AddInt32(&count, 1)
			}
			// for i := 0; i < n; i++ {
			// }
		}
	})
	assert(b, sl.Len() == int(count), "sl.GetSize() == count")
}

// Remove n elements per goroutines.
func BenchmarkRandomRemove(b *testing.B) {
	var sl = NewLockFreeSkipList(fn.Comp[int])
	var count int32
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			b.StopTimer()
			for i := 0; i < n; i++ {
				if sl.Add(rand.Int() % n) {
					atomic.AddInt32(&count, 1)
				}
			}
			b.StartTimer()
			for i := 0; i < n; i++ {
				if sl.Del(rand.Int() % n) {
					atomic.AddInt32(&count, -1)
				}
			}
		}
	})
	assert(b, sl.Len() == int(count), "sl.GetSize() == count")
}

func BenchmarkRandomGet(b *testing.B) {
	num := 102400
	var sl = NewLockFreeSkipList(fn.Comp[int])
	rand.Seed(time.Now().UnixNano())
	nums := []int{}
	for i := 0; i < num; i++ {
		nums = append(nums, rand.Intn(num))
	}
	for i := 0; i < num; i++ {
		// nums = append(nums, rand.Intn(num))
		sl.Add(nums[i])
	}
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// rand.Intn(num)
			for i := 0; i < num; i++ {
				sl.Get(nums[i])
			}
		}
	})

}

type SafeTreeMap[K any, V any] struct {
	tree  *LockFreeSkipList[K]
	mutex sync.Mutex
}

func (st *SafeTreeMap[K, V]) Put(key K, val V) bool {
	st.mutex.Lock()
	defer st.mutex.Unlock()
	return st.tree.Add(key)
}

func (st *SafeTreeMap[K, V]) Get(key K) K {
	st.mutex.Lock()
	defer st.mutex.Unlock()
	return st.tree.Get(key)
}

func BenchmarkRandomGetList(b *testing.B) {
	num := 1024
	var sl = &SafeTreeMap[int, int]{
		tree: NewLockFreeSkipList(fn.Comp[int]),
	}
	rand.Seed(time.Now().UnixNano())
	nums := []int{}
	for i := 0; i < num; i++ {
		nums = append(nums, rand.Intn(num))
	}
	for i := 0; i < num; i++ {
		// nums = append(nums, rand.Intn(num))
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

// Add and Remove n elements per goroutines.
func BenchmarkRandomAddAndRemoveAndContains(b *testing.B) {
	var sl = NewLockFreeSkipList(fn.Comp[int])
	divide := n / 3 // Make sure the total number of operations is n.
	var count int32
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			for i := 0; i < divide; i++ {
				if sl.Add(rand.Int() % divide) {
					atomic.AddInt32(&count, 1)
				}
			}
		}
	})
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			for i := 0; i < divide; i++ {
				if sl.Del(rand.Int() % divide) {
					atomic.AddInt32(&count, -1)
				}
			}
		}
	})
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			for i := 0; i < divide; i++ {
				sl.Has(rand.Int() % divide)
			}
		}
	})
	assert(b, sl.Len() == int(count), "sl.GetSize() == count")
}
