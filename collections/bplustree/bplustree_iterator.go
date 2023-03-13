package bplustree

// type Iterator[K, V any] struct {
// 	tree    *Tree[K, V]
// 	node    *Node[K, V]
// 	ptr     *kv.Pair[K, V]
// 	iterPos int32
// 	elemPos int32
// }

// func (iter *Iterator[K, V]) IsValid() bool {
// 	return iter.node != nil
// }

// func (iter *Iterator[K, V]) GetValue() kv.Entry[K, V] {
// 	return iter.ptr
// }

// func (iter *Iterator[K, V]) SetValue(p kv.Entry[K, V]) {
// 	iter.ptr.SetValue(p.GetValue())
// }

// const (
// 	end int32 = iota
// 	between
// )

// func (iter *Iterator[K, V]) becomeEnd() {
// 	iter.node = nil
// 	iter.ptr = nil
// 	iter.elemPos = -1
// 	iter.iterPos = end
// }

// func (iter *Iterator[K, V]) Next() {
// 	if iter.iterPos == end {
// 		left := iter.tree.LeftMost()
// 		if left == nil {
// 			iter.becomeEnd()
// 			return
// 		}
// 		iter.node = left
// 		iter.ptr = &left.values[0]
// 		iter.elemPos = 0
// 		iter.iterPos = between
// 		return
// 	}
// 	if int(iter.elemPos)+1 < iter.node.num {
// 		iter.elemPos++
// 		iter.ptr = &iter.node.values[iter.elemPos]
// 		return
// 	}
// 	left := iter.tree.LeftMost()
// 	if left == nil {
// 		iter.becomeEnd()
// 		return
// 	}
// 	if iter.node.next != nil && iter.node.next != left {
// 		iter.node = iter.node.next
// 		iter.ptr = &iter.node.values[0]
// 		iter.elemPos = 0
// 		return
// 	} else {
// 		iter.becomeEnd()
// 		return
// 	}
// }

// func (iter *Iterator[K, V]) Prev() {
// 	if iter.iterPos == end {
// 		right := iter.tree.RightMost()
// 		if right == nil {
// 			iter.becomeEnd()
// 			return
// 		}
// 		iter.node = right
// 		iter.ptr = &right.values[right.num-1]
// 		iter.elemPos = int32(right.num) - 1
// 		iter.iterPos = between
// 		return
// 	}
// 	if iter.elemPos-1 >= 0 {
// 		iter.ptr = &iter.node.values[iter.elemPos-1]
// 		iter.elemPos--
// 		return
// 	}
// 	right := iter.tree.RightMost()
// 	if right == nil {
// 		iter.node = nil
// 		iter.ptr = nil
// 		iter.elemPos = -1
// 		return
// 	}
// 	if iter.node.prev != nil && iter.node.prev != right {
// 		p := iter.node.prev
// 		iter.elemPos = int32(p.num - 1)
// 		iter.ptr = &iter.node.prev.values[p.num-1]
// 		iter.node = iter.node.prev
// 		return
// 	} else {
// 		iter.becomeEnd()
// 		return
// 	}

// }

// func (iter *Iterator[K, V]) Move(i int) {
// 	if i > 0 {
// 		for i != 0 && iter.IsValid() {
// 			iter.Next()
// 			i--
// 		}
// 	} else {
// 		for i != 0 && iter.IsValid() {
// 			iter.Prev()
// 			i++
// 		}
// 	}
// }

// func (iter *Iterator[K, V]) CompValue(p kv.Entry[K, V]) bool {
// 	return comparison.CompareWith(p.GetKey(), iter.ptr.Key, iter.tree.compFunc) == 0
// }

// func (iter *Iterator[K, V]) CompKey(key K) bool {
// 	return comparison.CompareWith(key, iter.ptr.Key, iter.tree.compFunc) == 0
// }

// func (iter *Iterator[K, V]) Comp(i iterator.Base[kv.Entry[K, V]]) bool {
// 	a, b := iter.IsValid(), i.IsValid()
// 	if a == b {
// 		if a {
// 			return i.(*Iterator[K, V]).ptr == iter.ptr
// 		}
// 		return true
// 	}
// 	return false
// }
// func (iter *Iterator[K, V]) CloneForward() iterator.Forward[kv.Entry[K, V]] {
// 	return &Iterator[K, V]{
// 		tree:    iter.tree,
// 		node:    iter.node,
// 		ptr:     iter.ptr,
// 		iterPos: iter.iterPos,
// 		elemPos: iter.elemPos,
// 	}
// }

// func (iter *Iterator[K, V]) CloneBidirectional() iterator.Bidirectional[kv.Entry[K, V]] {
// 	return &Iterator[K, V]{
// 		tree:    iter.tree,
// 		node:    iter.node,
// 		ptr:     iter.ptr,
// 		iterPos: iter.iterPos,
// 		elemPos: iter.elemPos,
// 	}
// }

// type ReverseIteratorAdapter[K, V any] struct {
// 	Iterator[K, V]
// }

// func (iter *ReverseIteratorAdapter[K, V]) Next() {
// 	iter.Prev()
// }

// func (iter *ReverseIteratorAdapter[K, V]) CloneForward() iterator.Forward[kv.Entry[K, V]] {
// 	return &ReverseIteratorAdapter[K, V]{iter.Iterator}
// }

// func InitUnaryIterator[T any](tree *Tree[T, constrains.Void]) UnaryIterator[T] {
// 	return UnaryIterator[T]{Iterator[T, constrains.Void]{tree: tree}}
// }

// func InitReverseUnaryIterator[T any](tree *Tree[T, constrains.Void]) UnaryReverseIteratorAdapter[T] {
// 	return UnaryReverseIteratorAdapter[T]{UnaryIterator[T]{Iterator[T, constrains.Void]{tree: tree}}}
// }

// type UnaryIterator[T any] struct {
// 	Iterator[T, constrains.Void]
// }

// func (iter *UnaryIterator[T]) GetValue() T {
// 	return iter.Iterator.GetValue().GetKey()
// }

// func (iter *UnaryIterator[T]) SetValue(p T) {}

// func (iter *UnaryIterator[T]) IsValid() bool {
// 	return iter.Iterator.IsValid()
// }

// func (iter *UnaryIterator[T]) Next() {
// 	iter.Iterator.IsValid()
// }

// func (iter *UnaryIterator[T]) Prev() {
// 	iter.Iterator.Prev()
// }

// func (iter *UnaryIterator[T]) Move(i int) {
// 	iter.Iterator.Move(i)
// }

// func (iter *UnaryIterator[T]) CompValue(k T) bool {
// 	return iter.Iterator.CompKey(k)
// }

// func (iter *UnaryIterator[T]) Comp(i iterator.Base[T]) bool {
// 	a, b := iter.IsValid(), i.IsValid()
// 	if a == b {
// 		if a {
// 			return i.(*UnaryIterator[T]).ptr == iter.ptr
// 		}
// 		return true
// 	}
// 	return false
// }

// func (iter *UnaryIterator[T]) CloneForward() iterator.Forward[T] {
// 	return &UnaryIterator[T]{iter.Iterator}
// }

// func (iter *UnaryIterator[T]) CloneBidirectional() iterator.Bidirectional[T] {
// 	return &UnaryIterator[T]{iter.Iterator}
// }

// // UnaryReverseIteratorAdapter it is adapter for original iterator to implement reverse by rewrite `Next` method
// // don't use interface or interface constrains for anonymous filed(UnaryIterator[T]), it will case extra cost
// type UnaryReverseIteratorAdapter[T any] struct {
// 	UnaryIterator[T]
// }

// func (iter *UnaryReverseIteratorAdapter[T]) Next() {
// 	iter.Prev()
// }

// func (iter *UnaryReverseIteratorAdapter[T]) CloneForward() iterator.Forward[T] {
// 	return &UnaryReverseIteratorAdapter[T]{iter.UnaryIterator}
// }
