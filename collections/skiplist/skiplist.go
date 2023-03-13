// Package skiplist : adapt from redis
package skiplist

import (
	"bytes"
	"fmt"
	"math/rand"
	"time"

	"github.com/ydmxcz/gds/collections/truple"
	"github.com/ydmxcz/gds/fn"
	"github.com/ydmxcz/gds/iterator"
)

// 这种写法docdog会解析错误
// const MaxLevel = 32

const (
	MaxLevel = 32
)

type List[K, V any] struct {
	header *Node[K, V]
	tail   *Node[K, V]
	size   int
	level  int
	random *rand.Rand
	comp   fn.Compare[K]
}

func NewNode[K, V any](level int, key K, val V) *Node[K, V] {
	return &Node[K, V]{
		value:    val,
		key:      key,
		backward: nil,
		level:    make([]Level[K, V], level),
	}
}

func New[K, V any](comp fn.Compare[K]) *List[K, V] {
	l := &List[K, V]{}
	l.Init(comp)
	return l
}

func (sl *List[K, V]) Init(comp fn.Compare[K]) {
	var val V
	var key K
	sl.header = NewNode(MaxLevel, key, val)
	sl.tail = nil
	sl.level = 1
	sl.size = 0
	sl.random = rand.New(rand.NewSource(time.Now().Unix()))
	sl.comp = comp
}

func (sl *List[K, V]) randomLevel() int {
	var total uint64 = 1<<MaxLevel - 1
	k := sl.random.Uint64() % total
	var levelN uint64 = 1 << (MaxLevel - 1)
	level := 1
	for total -= levelN; total > k; level++ {
		levelN >>= 1
		total -= levelN
	}
	return level
}

func (sl *List[K, V]) Put(key K, val V) bool {
	return sl.InsertUnique(key, val)
}

func (sl *List[K, V]) Has(key K) bool {
	n, _ := sl.search(key)
	if n == nil {
		return false
	}
	return true
}

func (sl *List[K, V]) Get(key K) (val V) {
	n, _ := sl.search(key)
	if n == nil {
		return
	}
	return n.value
}

func (sl *List[K, V]) search(key K) (*Node[K, V], int) {
	var x *Node[K, V]
	var rank int = 0
	var i int
	x = sl.header
	for i = sl.level - 1; i >= 0; i-- {
		//for x.level[i].forward != nil && comparison.CompareWith(x.level[i].forward.key, key, sl.comp) != 1 {
		for x.level[i].forward != nil && sl.comp(x.level[i].forward.key, key) != 1 {
			rank += x.level[i].span
			x = x.level[i].forward
		}
		if sl.comp(x.key, key) == 0 {
			return x, rank
		}
	}
	return nil, 0
}

func (sl *List[K, V]) Delete(key K) bool {
	var update [MaxLevel]*Node[K, V]
	var x *Node[K, V]
	var i int
	x = sl.header
	for i = sl.level - 1; i >= 0; i-- {
		for x.level[i].forward != nil && sl.comp(x.level[i].forward.key, key) == 0 {
			x = x.level[i].forward
		}
		update[i] = x
	}
	x = x.level[0].forward
	if x != nil && sl.comp(key, x.key) == 0 {
		// delete node
		sl.deleteNode(x, &update)
		x = nil
		return true
	}
	return false
}

func (sl *List[K, V]) DeleteByWholeMapping(key K, val V, valComp fn.Compare[V]) bool {
	n, _ := sl.search(key)
	if n != nil && valComp(val, n.value) == 0 {
		sl.Delete(key)
		return true
	}
	return false
}

func (sl *List[K, V]) DeleteKeys(keys ...K) int {
	changed := 0
	for i := 0; i < len(keys); i++ {
		if sl.Delete(keys[i]) {
			changed++
		}
	}
	return changed
}

func (sl *List[K, V]) Clear() {
	comp := sl.comp
	sl.Init(comp)
}

func (sl *List[K, V]) IsEmpty() bool {
	return sl.size == 0
}

func (sl *List[K, V]) String() string {
	var buf bytes.Buffer
	var next *Node[K, V]
	buf.WriteString("[")
	for n := sl.First(); n.IsValid(); {
		buf.WriteString(fmt.Sprintf("%v:%v", n.key, n.value))
		next = n.Next()
		if next != nil {
			buf.WriteByte(' ')
		}
		n = next
	}
	buf.WriteByte(']')
	return buf.String()
}

func (sl *List[K, V]) SetCompareFunc(compFunc fn.Compare[K]) {
	if sl.comp == nil {
		sl.comp = compFunc
	}
}

func (sl *List[K, V]) GetCompareFunc() fn.Compare[K] {
	return sl.comp
}

func (sl *List[K, V]) First() *Node[K, V] {
	return sl.header.level[0].forward
}

func (sl *List[K, V]) Last() *Node[K, V] {
	return sl.tail
}

func (sl *List[K, V]) Size() int {
	return int(sl.size)
}

func (sl *List[K, V]) InsertUnique(key K, val V) bool {
	return sl.insert(key, val, true)
}

func (sl *List[K, V]) InsertEquals(key K, val V) bool {
	return sl.insert(key, val, false)
}

func (sl *List[K, V]) Iter() iterator.Iter[truple.KV[K, V]] {
	return sl.iter(sl.First())
}

func (sl *List[K, V]) iter(node *Node[K, V]) iterator.Iter[truple.KV[K, V]] {
	return func() (val truple.KV[K, V], ok bool) {
		if ok = node != nil; ok {
			val = truple.KV[K, V]{Val: node.value, Key: node.key}
			node = node.Next()
		}
		return
	}
}

func (sl *List[K, V]) iterStep(step int, node *Node[K, V]) iterator.Iter[truple.KV[K, V]] {
	return func() (val truple.KV[K, V], ok bool) {
		if ok = (step > 0 && node != nil); ok {
			val = truple.KV[K, V]{Val: node.value, Key: node.key}
			node = node.Next()
		}
		step--
		return
	}
}

func (sl *List[K, V]) SplitableIter() func(parallelism int) iterator.Iter[iterator.Iter[truple.KV[K, V]]] {
	return func(parallelism int) iterator.Iter[iterator.Iter[truple.KV[K, V]]] {
		idx := 0
		var step int
		if parallelism == 0 {
			step = int(sl.size)
		} else {
			step = int(sl.size) / parallelism
		}
		if parallelism <= 0 {
			return func() (iterator.Iter[truple.KV[K, V]], bool) {
				if idx == 0 {
					idx++
					return sl.Iter(), true
				}
				return nil, false
			}
		}
		node := sl.First()
		return func() (iterator.Iter[truple.KV[K, V]], bool) {
			if idx >= int(sl.size) {
				return nil, false
			}
			i := idx
			idx += step

			n := node
			for j := i; j < i+step && node != nil; j++ {
				node = node.Next()
			}

			if i+step >= sl.size {
				return sl.iter(n), true
			}
			return sl.iterStep(step, n), true
		}
	}
}

func (sl *List[K, V]) insert(key K, val V, isUnique bool) bool {
	var update [MaxLevel]*Node[K, V]
	var x *Node[K, V]
	var rank [MaxLevel]int
	var i, level int
	x = sl.header
	for i = sl.level - 1; i >= 0; i-- {
		/* store rank that is crossed to reach the insert position */
		if i == (sl.level - 1) {
			rank[i] = 0
		} else {
			rank[i] = rank[i+1]
		}
		for x.level[i].forward != nil && sl.comp(x.level[i].forward.key, key) == 0 {
			rank[i] += x.level[i].span
			x = x.level[i].forward
		}
		update[i] = x
	}

	if isUnique && update[0].level[0].forward != nil && sl.comp(update[0].level[0].forward.key, key) == 0 {
		return false
	}

	/* we assume the value is not already inside, since we allow duplicated
	 * scores, reinserting the same value should never happen since the
	 * caller of zslInsert() should test in the hash table if the value is
	 * already inside or not. */
	level = sl.randomLevel()
	if level > sl.level {
		for i = sl.level; i < level; i++ {
			rank[i] = 0
			update[i] = sl.header
			update[i].level[i].span = sl.size
		}
		sl.level = level
	}
	x = NewNode(level, key, val)
	for i = 0; i < level; i++ {
		x.level[i].forward = update[i].level[i].forward
		update[i].level[i].forward = x

		/* update span covered by update[i] as x is inserted here */
		x.level[i].span = update[i].level[i].span - (rank[0] - rank[i])
		update[i].level[i].span = (rank[0] - rank[i]) + 1
	}

	/* increment span for untouched levels */
	for i = level; i < sl.level; i++ {
		update[i].level[i].span++
	}
	if update[0] == sl.header {
		x.backward = nil
	} else {
		x.backward = update[0]
	}
	// x.backward = () ? NULL : update[0];
	if x.level[0].forward != nil {
		x.level[0].forward.backward = x

	} else {
		sl.tail = x
	}
	sl.size++
	return true
}

func (sl *List[K, V]) deleteNode(x *Node[K, V], update *[MaxLevel]*Node[K, V]) {
	var i int
	for i = 0; i < sl.level; i++ {
		if update[i].level[i].forward == x {
			update[i].level[i].span += x.level[i].span - 1
			update[i].level[i].forward = x.level[i].forward
		} else {
			update[i].level[i].span -= 1
		}
	}
	if x.level[0].forward != nil {
		x.level[0].forward.backward = x.backward
	} else {
		sl.tail = x.backward
	}
	for sl.level > 1 && sl.header.level[sl.level-1].forward == nil {
		sl.level--
	}
	sl.size--
}

func (sl *List[K, V]) Rank(key K) int {
	n, rank := sl.search(key)
	if n != nil {
		return rank
	}
	return 0
}

func (sl *List[K, V]) GetByRank(rank int) (val V) {
	var x *Node[K, V]
	var traversed int
	var i int
	x = sl.header
	for i = sl.level - 1; i >= 0; i-- {
		for x.level[i].forward != nil && (traversed+x.level[i].span) <= rank {
			traversed += x.level[i].span
			x = x.level[i].forward
		}
		if traversed == rank {
			return x.value
		}
	}
	return val
}

// IsInAllRange Returns if there is a part of the zset is in range.
// this util includes the maximum value of the skip linkedlist
func (sl *List[K, V]) IsInAllRange(min, max K) bool {
	var x *Node[K, V]
	if sl.comp(min, max) > -1 {
		// if min >= max {
		return false
	}
	x = sl.tail
	// The max value of skip linkedlist less than the min value,
	if x == nil || sl.comp(x.key, min) == -1 || sl.comp(x.key, max) == -1 {
		return false
	}
	// The min value of skip linkedlist greater than the max value.
	x = sl.header.level[0].forward
	if x == nil || sl.comp(x.key, max) == 1 || sl.comp(x.key, min) == 1 {
		return false
	}
	return true
}

// IsInPartOfRange Returns if there is a part of the zset is in range.
// this util includes the maximum value of the skip linkedlist
func (sl *List[K, V]) IsInPartOfRange(min, max K) bool {
	// fmt.Println(sl.level)
	var x *Node[K, V]
	if sl.comp(min, max) > -1 {
		return false
	}
	x = sl.tail
	if x == nil || sl.comp(x.key, min) == -1 {
		return false
	}
	x = sl.header.level[0].forward
	if x == nil || sl.comp(x.key, max) == 1 {
		return false
	}
	return true
}

func (sl *List[K, V]) FirstInRange(min, max K) *Node[K, V] {
	if !sl.IsInAllRange(min, max) {
		return nil
	}
	// var x *Node[K,V]
	x := sl.header
	for i := sl.level - 1; i >= 0; i-- {
		for x.level[i].forward != nil && sl.comp(x.level[i].forward.key, min) == -1 {
			x = x.level[i].forward
		}
	}
	x = x.level[0].forward
	if x == nil {
		panic("x == nil")
	}
	if sl.comp(x.key, max) == 1 {
		return nil
	}
	return x
}

func (sl *List[K, V]) LastInRange(min, max K) *Node[K, V] {
	if !sl.IsInAllRange(min, max) {
		return nil
	}
	x := sl.header
	for i := sl.level - 1; i >= 0; i-- {
		for x.level[i].forward != nil && sl.comp(x.level[i].forward.key, max) == -1 {
			x = x.level[i].forward
		}
	}
	x = x.level[0].forward
	if x == nil {
		panic("x == nil")
	}
	if sl.comp(x.key, min) == -1 {

		return nil
	}
	return x
}

func (sl *List[K, V]) DeleteRangeByKey(min, max K) uint64 {
	var update [MaxLevel]*Node[K, V]
	var x, next *Node[K, V]
	var removed uint64
	i := 0
	x = sl.header
	for i = sl.level - 1; i >= 0; i-- {
		for x.level[i].forward != nil && sl.comp(x.level[i].forward.key, min) == -1 {
			x = x.level[i].forward
		}
		update[i] = x
	}

	x = x.level[0].forward

	for x != nil && sl.comp(x.key, max) == -1 {
		next = x.level[0].forward
		sl.deleteNode(x, &update)
		removed++
		x = next
	}
	return removed
}
