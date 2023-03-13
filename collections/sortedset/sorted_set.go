package sortedset

import (
	"sync"

	"github.com/alphadose/haxmap"
	"github.com/ydmxcz/gds/collections/rbtree"
	"github.com/ydmxcz/gds/util/constraints"
)

type SortedSet[K constraints.Hashable, V any] struct {
	hashmap *haxmap.Map[K, V]
	mutex   sync.Mutex
	tree    *rbtree.Tree[V, K]
}

func New[K constraints.Hashable, V any]() *SortedSet[K, V] {
	// haxmap.New[]()
	return nil
}

func (set *SortedSet[K, V]) Set(key K, val V) {
	_, loaded := set.hashmap.GetOrSet(key, val)
	if loaded {
		set.mutex.Lock()
		n := set.tree.Delete(val)
		set.mutex.Unlock()
		n.SetKey(val)
		set.mutex.Lock()
		set.tree.PutNode(n)
		set.mutex.Unlock()
		return
	} else {
		set.tree.Put(val, key)
	}
}

func (set *SortedSet[K, V]) Get(key K) (val V, exist bool) {
	return set.hashmap.Get(key)
}

// func (set *SortedSet[K, V]) Rank() []K {

// }
