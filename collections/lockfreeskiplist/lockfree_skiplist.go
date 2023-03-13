package lockfreeskiplist

import (
	"math/rand"
	"sync/atomic"
	"unsafe"
)

const maxLevel = 20

// LockFreeSkipList define
type LockFreeSkipList[T any] struct {
	head *node[T]
	tail *node[T]
	size int32
	comp func(value1 T, value2 T) bool
}

type node[T any] struct {
	level int
	nexts []unsafe.Pointer
	value T
}

// NewLockFreeSkipList new a lockfree skiplist, you should pass a compare function.
func NewLockFreeSkipList[T any](comp func(value1 T, value2 T) bool) *LockFreeSkipList[T] {
	sl := new(LockFreeSkipList[T])
	var defVal T
	sl.head = newNode(maxLevel, defVal)
	sl.tail = newNode(maxLevel, defVal)
	sl.size = 0
	sl.comp = comp
	for level := 0; level < maxLevel; level++ {
		sl.head.storeNext(level, sl.tail)
	}
	return sl
}

// Add a value to skiplist.
func (sl *LockFreeSkipList[T]) Add(value T) bool {
	var prevs [maxLevel]*node[T]
	var nexts [maxLevel]*node[T]
	for true {
		if sl.find(value, &prevs, &nexts) {
			return false
		}
		topLevel := randomLevel()
		newNode := newNode(topLevel, value)
		for level := 0; level < topLevel; level++ {
			newNode.storeNext(level, nexts[level])
		}
		if prev, next := prevs[0], nexts[0]; !prev.casNext(0, next, newNode) {
			// The successor of prev is not next, we should try again.
			continue
		}
		for level := 1; level < topLevel; level++ {
			for true {
				if prev, next := prevs[level], nexts[level]; prev.casNext(level, next, newNode) {
					break
				}
				// The successor of prev is not next,
				// we should call find to update the prevs and nexts.
				sl.find(value, &prevs, &nexts)
			}
		}
		break
	}
	atomic.AddInt32(&sl.size, 1)
	return true
}

// Del a value from skiplist.
func (sl *LockFreeSkipList[T]) Del(value T) bool {
	var prevs [maxLevel]*node[T]
	var nexts [maxLevel]*node[T]
	if !sl.find(value, &prevs, &nexts) {
		return false
	}
	removeNode := nexts[0]
	for level := removeNode.level - 1; level > 0; level-- {
		next := removeNode.loadNext(level)
		for !isMarked(next) {
			// Make sure that all but the bottom next are marked from top to bottom.
			removeNode.casNext(level, next, getMarked(next))
			next = removeNode.loadNext(level)
		}
	}
	for next := removeNode.loadNext(0); true; next = removeNode.loadNext(0) {
		if isMarked(next) {
			// Other thread already maked the next, so this thread delete failed.
			return false
		}
		if removeNode.casNext(0, next, getMarked(next)) {
			// This thread marked the bottom next, delete successfully.
			break
		}
	}
	atomic.AddInt32(&sl.size, -1)
	return true
}

// Has check if skiplist contains a value.
func (sl *LockFreeSkipList[T]) Has(value T) bool {
	var prevs [maxLevel]*node[T]
	var nexts [maxLevel]*node[T]
	return sl.find(value, &prevs, &nexts)
}

func (sl *LockFreeSkipList[T]) Get(value T) T {
	var prevs [maxLevel]*node[T]
	var nexts [maxLevel]*node[T]
	return sl.get(value, &prevs, &nexts)
}

func (sl *LockFreeSkipList[T]) GetRank(value T) int {
	var prevs [maxLevel]*node[T]
	var nexts [maxLevel]*node[T]
	return sl.getRank(value, &prevs, &nexts)
}

// GetSize get the element size of skiplist.
func (sl *LockFreeSkipList[T]) GetSize() int32 {
	return atomic.LoadInt32(&sl.size)
}

func randomLevel() int {
	level := 1
	for level < maxLevel && rand.Int()&1 == 0 {
		level++
	}
	return level
}

func (sl *LockFreeSkipList[T]) less(nd *node[T], value T) bool {
	if sl.head == nd {
		return true
	}
	if sl.tail == nd {
		return false
	}
	return sl.comp(nd.value, value)
}

func (sl *LockFreeSkipList[T]) equals(nd *node[T], value T) bool {
	if sl.head == nd || sl.tail == nd {
		return false
	}
	return !sl.comp(nd.value, value) && !sl.comp(value, nd.value)
}

func (sl *LockFreeSkipList[T]) get(value T, prevs *[maxLevel]*node[T], nexts *[maxLevel]*node[T]) T {
	var prev *node[T]
	var cur *node[T]
	var next *node[T]
retry:
	prev = sl.head
	for level := maxLevel - 1; level >= 0; level-- {
		cur = getUnmarked(prev.loadNext(level))
		for true {
			next = cur.loadNext(level)
			for isMarked(next) {
				// Like harris-linkedlist,remove the node while traversing.
				// See also https://github.com/bhhbazinga/LockFreeLinkedList.
				if !prev.casNext(level, cur, getUnmarked(next)) {
					goto retry
				}
				cur = getUnmarked(prev.loadNext(level))
				next = cur.loadNext(level)
			}
			if !sl.less(cur, value) {
				break
			}
			prev = cur
			cur = next
		}
		prevs[level] = prev
		nexts[level] = cur
	}
	return cur.value
	// return sl.equals(cur, value)
}

func (sl *LockFreeSkipList[T]) getRank(value T, prevs *[maxLevel]*node[T], nexts *[maxLevel]*node[T]) int {
	var prev *node[T]
	var cur *node[T]
	var next *node[T]
retry:
	prev = sl.head
	for level := maxLevel - 1; level >= 0; level-- {
		cur = getUnmarked(prev.loadNext(level))
		for true {
			next = cur.loadNext(level)
			for isMarked(next) {
				// Like harris-linkedlist,remove the node while traversing.
				// See also https://github.com/bhhbazinga/LockFreeLinkedList.
				if !prev.casNext(level, cur, getUnmarked(next)) {
					goto retry
				}
				cur = getUnmarked(prev.loadNext(level))
				next = cur.loadNext(level)
			}
			if !sl.less(cur, value) {
				break
			}
			prev = cur
			cur = next
		}
		prevs[level] = prev
		nexts[level] = cur
	}
	return -1
	// return sl.equals(cur, value)
}

func (sl *LockFreeSkipList[T]) find(value T, prevs *[maxLevel]*node[T], nexts *[maxLevel]*node[T]) bool {
	var prev *node[T]
	var cur *node[T]
	var next *node[T]
retry:
	prev = sl.head
	for level := maxLevel - 1; level >= 0; level-- {
		cur = getUnmarked(prev.loadNext(level))
		for true {
			next = cur.loadNext(level)
			for isMarked(next) {
				// Like harris-linkedlist,remove the node while traversing.
				// See also https://github.com/bhhbazinga/LockFreeLinkedList.
				if !prev.casNext(level, cur, getUnmarked(next)) {
					goto retry
				}
				cur = getUnmarked(prev.loadNext(level))
				next = cur.loadNext(level)
			}
			if !sl.less(cur, value) {
				break
			}
			prev = cur
			cur = next
		}
		prevs[level] = prev
		nexts[level] = cur
	}
	return sl.equals(cur, value)
}

func newNode[T any](level int, value T) *node[T] {
	nd := new(node[T])
	nd.nexts = make([]unsafe.Pointer, level)
	nd.level = level
	nd.value = value
	for level := 0; level < nd.level; level++ {
		nd.storeNext(level, nil)
	}
	return nd
}

func (nd *node[T]) loadNext(level int) *node[T] {
	return (*node[T])(atomic.LoadPointer(&nd.nexts[level]))
}

func (nd *node[T]) storeNext(level int, next *node[T]) {
	atomic.StorePointer(&nd.nexts[level], unsafe.Pointer(next))
}

func (nd *node[T]) casNext(level int, expected *node[T], desire *node[T]) bool {
	return atomic.CompareAndSwapPointer(&nd.nexts[level], unsafe.Pointer(expected), unsafe.Pointer(desire))
}

func isMarked[T any](next *node[T]) bool {
	ptr := unsafe.Pointer(next)
	return uintptr(ptr)&uintptr(0x1) == 1
}

func getMarked[T any](next *node[T]) *node[T] {
	ptr := unsafe.Pointer(next)
	return (*node[T])(unsafe.Pointer(uintptr(ptr) | uintptr(0x1)))

}

func getUnmarked[T any](next *node[T]) *node[T] {
	ptr := unsafe.Pointer(next)
	return (*node[T])(unsafe.Pointer(uintptr(ptr) & ^uintptr(0x1)))
}
