// Package bplustree adapt from: https://github.com/orange1438/BTree-and-BPlusTree-Realize
package bplustree

import (
	"fmt"

	// "github.com/edmundshelby/genlib/util/comparison"

	// "github.com/edmundshelby/genlib/util/comparison"

	"github.com/ydmxcz/gds/collections/truple"
	"github.com/ydmxcz/gds/fn"
	// "github.com/edmundshelby/genlib/collections/associative/kv"
	// "github.com/edmundshelby/genlib/iterator"
)

type Tree[K, V any] struct {
	root     *Node[K, V]
	size     int
	degree   int
	compFunc fn.Compare[K]
}

func NewMap[K, V any](degree int, compFunc fn.Compare[K]) *Tree[K, V] {
	t := &Tree[K, V]{}
	t.Init(degree, compFunc)
	return t
}

func (bt *Tree[K, V]) Init(degree int, compFunc fn.Compare[K]) {
	bt.root = NewNode[K, V](degree)
	bt.size = 0
	bt.degree = degree
	bt.compFunc = compFunc
	bt.root.prev = bt.root
	bt.root.next = bt.root
}

func NewNode[K, V any](degree int) *Node[K, V] {
	n := &Node[K, V]{
		values:   make([]truple.KV[K, V], (degree*2)-1),
		children: make([]*Node[K, V], degree*2),
		num:      0,
		isLeaf:   true,
	}
	return n
}

func (bt *Tree[K, V]) ContainsKey(key K) bool {
	_, _, exist := bt.search(bt.root, key)
	return exist
}

func (bt *Tree[K, V]) Get(key K) (val V) {
	node, idx, exist := bt.search(bt.root, key)
	if exist {
		return node.values[idx].Val
	}
	return
}

func (bt *Tree[K, V]) GetEntry(key K) *truple.KV[K, V] {
	node, idx, exist := bt.search(bt.root, key)
	if exist {
		return &node.values[idx]
	}
	return nil
}

func (bt *Tree[K, V]) GetVerify(key K) (val V, exist bool) {
	node, idx, exist := bt.search(bt.root, key)
	if exist {
		return node.values[idx].Val, true
	}
	return
}

func (bt *Tree[K, V]) DeleteByWholeMapping(key K, val V, compV fn.Compare[V]) bool {
	node, idx, exist := bt.search(bt.root, key)
	if exist && compV(node.values[idx].Val, val) < 0 {
		return bt.Delete(key)
	}
	return false
}

func (bt *Tree[K, V]) DeleteByKeys(keys ...K) int {
	changed := 0
	for i := 0; i < len(keys); i++ {
		if bt.Delete(keys[i]) {
			changed++
		}
	}
	return changed
}

func (bt *Tree[K, V]) DeleteEntryByKey(entry truple.KV[K, V]) bool {
	return bt.Delete(entry.Key)
}

func (bt *Tree[K, V]) Delete(key K) bool {
	if bt.ContainsKey(key) {
		var p truple.KV[K, V]
		p.Key = key
		bt.root = bt.delete(bt.root, &p)
		//bt.root.parent = nil
		return true
	}
	return false
}

func (bt *Tree[K, V]) splitChild(parent, child *Node[K, V], pos int) {
	M := bt.degree
	newChild := NewNode[K, V](M)

	newChild.isLeaf = child.isLeaf
	newChild.num = M - 1

	for i := 0; i < M-1; i++ {
		newChild.values[i] = child.values[i+M]
	}

	if false == newChild.isLeaf {
		for i := 0; i < M; i++ {
			newChild.children[i] = child.children[i+M]
		}
	}

	child.num = M - 1
	if true == child.isLeaf {
		child.num++ // if is leaf, keep the middle ele, put it in the left
	}

	for i := parent.num; i > pos; i-- {
		parent.children[i+1] = parent.children[i]
	}
	parent.children[pos+1] = newChild

	for i := parent.num - 1; i >= pos; i-- {
		parent.values[i+1] = parent.values[i]
	}
	parent.values[pos] = child.values[M-1]

	parent.num += 1

	// update link
	if true == child.isLeaf {
		newChild.next = child.next
		child.next.prev = newChild
		newChild.prev = child
		child.next = newChild
	}
}

func (bt *Tree[K, V]) insertNonFull(node *Node[K, V], p *truple.KV[K, V]) {
	if true == node.isLeaf {
		pos := node.num
		for pos >= 1 && bt.compFunc(p.Key, node.values[pos-1].Key) < 0 {
			node.values[pos] = node.values[pos-1]
			pos--
		}

		node.values[pos] = *p
		node.num += 1
		bt.size += 1

	} else {
		pos := node.num
		for pos > 0 && bt.compFunc(p.Key, node.values[pos-1].Key) < 0 {
			pos--
		}

		if 2*bt.degree-1 == node.children[pos].num {
			bt.splitChild(node, node.children[pos], pos)
			if bt.compFunc(node.values[pos].Key, p.Key) < 0 {
				pos++
			}
		}

		bt.insertNonFull(node.children[pos], p)
	}
}

func (bt *Tree[K, V]) insert(root *Node[K, V], p *truple.KV[K, V]) *Node[K, V] {
	if root == nil {
		return nil
	}
	if (bt.degree*2)-1 == root.num {
		node := NewNode[K, V](bt.degree)
		node.isLeaf = false
		node.children[0] = root
		bt.splitChild(node, root, 0)
		bt.insertNonFull(node, p)
		return node
	} else {
		bt.insertNonFull(root, p)
		return root
	}
}

func (bt *Tree[K, V]) search(startNode *Node[K, V], key K) (node *Node[K, V], index int, exist bool) {
	if bt.size == 0 {
		return nil, -1, false
	}
	node = startNode
	var l, r, mid int
	var p *truple.KV[K, V]
	for {
		exist = false
		if node != nil {
			l, r, mid = 0, node.num-1, 0
			for l <= r {
				mid = (l + r) >> 1
				p = &node.values[mid]
				compare := bt.compFunc(key, p.Key)
				if compare < 0 {
					r = mid - 1
				} else if compare > 0 {
					l = mid + 1
				} else {
					exist = true
					break
				}
			}
		}
		if node.isLeaf {
			if exist {
				return node, mid, true
			}
			return node, l, false
		}
		node = node.children[l]
	}
}

func (bt *Tree[K, V]) InsertEqual(key K, val V) bool {
	var p truple.KV[K, V]
	p.Key = key
	p.Val = val
	bt.root = bt.insert(bt.root, &p)
	return true
}

func (bt *Tree[K, V]) InsertUnique(key K, val V) bool {
	if !bt.ContainsKey(key) {
		var p truple.KV[K, V]
		p.Key = key
		p.Val = val
		bt.root = bt.insert(bt.root, &p)
		return true
	}
	return false
}

func (bt *Tree[K, V]) mergeChild(root, y, z *Node[K, V], pos int) {
	M := bt.degree
	if true == y.isLeaf {
		y.num = 2*M - 2
		for i := M; i < 2*M-1; i++ {
			y.values[i-1] = z.values[i-M]
		}
	} else {
		y.num = 2*M - 1
		for i := M; i < 2*M-1; i++ {
			y.values[i] = z.values[i-M]
		}
		y.values[M-1] = root.values[pos]
		for i := M; i < 2*M; i++ {
			y.children[i] = z.children[i-M]
		}
	}

	for j := pos + 1; j < root.num; j++ {
		root.values[j-1] = root.values[j]
		root.children[j] = root.children[j+1]
	}

	root.num -= 1

	// update link
	if true == y.isLeaf {
		y.next = z.next
		z.next.prev = y
	}
	z = nil
	//free(z);
}

// Size return the value which equals key
func (bt *Tree[K, V]) Size() int {
	return bt.size
}

// IsEmpty return the value which equals key
func (bt *Tree[K, V]) IsEmpty() bool {
	return bt.size == 0
}

// Clear return the value which equals key
func (bt *Tree[K, V]) Clear() {
	bt.size = 0
	bt.root = nil
}

func (bt *Tree[K, V]) LeftMost() *Node[K, V] {
	return bt.root.MinimumNode()
}

func (bt *Tree[K, V]) RightMost() *Node[K, V] {
	return bt.root.MaximumNode()
}

func (bt *Tree[K, V]) DeleteAll(keys []K) int {
	changed := 0
	for i := 0; i < len(keys); i++ {
		if bt.Delete(keys[i]) {
			changed++
		}
	}
	return changed
}

func (bt *Tree[K, V]) delete(root *Node[K, V], p *truple.KV[K, V]) *Node[K, V] {
	if 1 == root.num {
		y := root.children[0]
		z := root.children[1]
		if nil != y && nil != z &&
			bt.degree-1 == y.num && bt.degree-1 == z.num {
			bt.mergeChild(root, y, z, 0)
			//free(root);
			root = nil
			bt.deleteNonOne(y, p)
			return y
		} else {
			bt.deleteNonOne(root, p)
			return root
		}
	} else {
		bt.deleteNonOne(root, p)
		return root
	}
}

func (bt *Tree[K, V]) deleteNonOne(root *Node[K, V], p *truple.KV[K, V]) {
	M := bt.degree
	if true == root.isLeaf {
		i := 0
		for i < root.num && bt.compFunc(root.values[i].Key, p.Key) < 0 {
			i++
		}
		if bt.compFunc(p.Key, root.values[i].Key) == 0 {
			for j := i + 1; j < 2*M-1; j++ {
				root.values[j-1] = root.values[j]
			}
			root.num -= 1
			bt.size -= 1

		} else {
			return
		}
	} else {
		i := 0
		var y, z *Node[K, V]
		for i < root.num && bt.compFunc(root.values[i].Key, p.Key) < 0 {
			i++
		}

		y = root.children[i]
		if i < root.num {
			z = root.children[i+1]
		}
		var t *Node[K, V]
		if i > 0 {
			t = root.children[i-1]
		}

		if y.num == M-1 {
			if i > 0 && t.num > M-1 {
				bt.shiftToRightChild(root, t, y, i-1)
			} else if i < root.num && z.num > M-1 {
				bt.shiftToLeftChild(root, y, z, i)
			} else if i > 0 {
				bt.mergeChild(root, t, y, i-1)
				y = t
			} else {
				bt.mergeChild(root, y, z, i)
			}
			bt.deleteNonOne(y, p)
		} else {
			bt.deleteNonOne(y, p)
		}
	}
}

func (bt *Tree[K, V]) searchPreSuccessor(root *Node[K, V]) *truple.KV[K, V] {
	y := root
	for false == y.isLeaf {
		y = y.children[y.num]
	}
	return &y.values[y.num-1]
}

func (bt *Tree[K, V]) searchSuccessor(root *Node[K, V]) *truple.KV[K, V] {
	z := root
	for false == z.isLeaf {
		z = z.children[0]
	}
	return &z.values[0]
}

func (bt *Tree[K, V]) shiftToRightChild(root, y, z *Node[K, V], pos int) {
	z.num += 1

	if false == z.isLeaf {
		z.values[0] = root.values[pos]
		root.values[pos] = y.values[y.num-1]
	} else {
		z.values[0] = y.values[y.num-1]
		root.values[pos] = y.values[y.num-2]
	}

	for i := z.num - 1; i > 0; i-- {
		z.values[i] = z.values[i-1]
	}

	if false == z.isLeaf {
		for i := z.num; i > 0; i-- {
			z.children[i] = z.children[i-1]
		}
		z.children[0] = y.children[y.num]
	}

	y.num -= 1
}

func (bt *Tree[K, V]) shiftToLeftChild(root, y, z *Node[K, V], pos int) {
	y.num += 1

	if false == z.isLeaf {
		y.values[y.num-1] = root.values[pos]
		root.values[pos] = z.values[0]
	} else {
		y.values[y.num-1] = z.values[0]
		root.values[pos] = z.values[0]
	}

	for j := 1; j < z.num; j++ {
		z.values[j-1] = z.values[j]
	}

	if false == z.isLeaf {
		y.children[y.num] = z.children[0]
		for j := 1; j <= z.num; j++ {
			z.children[j-1] = z.children[j]
		}
	}

	z.num -= 1
}

func (bt *Tree[K, V]) display() {
	var queue = make([]*Node[K, V], bt.size+1)
	front := 0
	rear := 0

	queue[rear] = bt.root
	rear++
	for front < rear {
		node := queue[front]
		front++
		fmt.Print("[")
		for i := 0; i < node.num; i++ {
			fmt.Printf("%v ", node.values[i].Key)
		}
		fmt.Printf("] ")

		for i := 0; i <= node.num; i++ {
			if nil != node.children[i] {
				queue[rear] = node.children[i]
				rear++
			}
		}
	}
	fmt.Println()
}

func (bt *Tree[K, V]) GetCompareFunc() fn.Compare[K] {
	return bt.compFunc
}
