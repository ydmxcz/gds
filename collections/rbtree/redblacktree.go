package rbtree

import (
	"bytes"
	"fmt"

	"github.com/ydmxcz/gds/collections/truple"
	"github.com/ydmxcz/gds/fn"
	"github.com/ydmxcz/gds/iterator"
	"github.com/ydmxcz/gds/stream"
	"github.com/ydmxcz/gds/util/constraints"
)

type Tree[K any, V any] struct {
	header     Node[K, V]
	comparator fn.Compare[K]
	size       int
	buf        nodeBuffer[K, V]
}

func NewOrderd[K constraints.Ordered, V any]() *Tree[K, V] {
	return New[K, V](fn.Comp[K])
}

func New[K any, V any](comp fn.Compare[K]) *Tree[K, V] {
	t := &Tree[K, V]{}
	t.Init(comp)
	return t
}

func (rbTree *Tree[K, V]) Init(lessFunc fn.Compare[K]) {
	rbTree.comparator = lessFunc
	rbTree.header.parent = nil
	rbTree.header.right = &rbTree.header
	rbTree.header.left = &rbTree.header
	rbTree.size = 0
	rbTree.buf.Init(16)
}

func (rbTree *Tree[K, V]) allocNode(key K, val V) *Node[K, V] {
	n := rbTree.buf.Alloc()
	n.kv = truple.KV[K, V]{
		Val: val,
		Key: key,
	}
	n.color = Red
	return n
}

func (rbTree *Tree[K, V]) freeNode(n *Node[K, V]) {
	rbTree.buf.Free(n)
}

func (rbTree *Tree[K, V]) String() string {
	var buf bytes.Buffer
	var next *Node[K, V]
	buf.WriteString("[")
	for n := rbTree.header.left; n.IsValid(); {
		buf.WriteString(fmt.Sprintf("%v:%v", n.kv.Key, n.kv.Val))
		next = n.Next()
		if next != nil {
			buf.WriteByte(' ')
		}
		n = next
	}
	buf.WriteByte(']')
	return buf.String()
}

func (rbTree *Tree[K, V]) SetCompareFunc(compFunc fn.Compare[K]) {
	if rbTree.comparator == nil {
		rbTree.comparator = compFunc
	}
}

func (rbTree *Tree[K, V]) GetCompareFunc() fn.Compare[K] {
	return rbTree.comparator
}

// left rotating are adjusts the position of x's right child and x.
func leftRotate[K, V any](header, x *Node[K, V]) {
	// y is the x's right child
	y := x.right
	// Changing x's right child and y's left child
	x.right = y.left
	if y.left != nil {
		y.left.parent = x
	}
	// move x's parent to y
	y.parent = x.parent
	// make x's parent's children is y
	if x.parent == nil {
		header.parent = y
	} else if x == x.parent.left {
		x.parent.left = y
	} else {
		x.parent.right = y
	}
	// change the parent relation of y and x
	y.left = x
	x.parent = y
}

func rightRotate[K, V any](header, x *Node[K, V]) {
	y := x.left

	x.left = y.right
	if y.right != nil {
		y.right.parent = x
	}

	y.parent = x.parent

	if x.parent == nil {
		header.parent = y
	} else if x == x.parent.right {
		x.parent.right = y
	} else {
		x.parent.left = y
	}

	y.right = x
	x.parent = y
}

func (rbTree *Tree[K, V]) RightMost() *Node[K, V] {
	return rbTree.header.right
}

func (rbTree *Tree[K, V]) LeftMost() *Node[K, V] {
	return rbTree.header.left
}

func (rbTree *Tree[K, V]) Clear() {
	rbTree.header.parent = nil
	rbTree.header.right = &rbTree.header
	rbTree.header.left = &rbTree.header
	rbTree.size = 0
}

func (rbTree *Tree[K, V]) IsEmpty() bool {
	return rbTree.size == 0
}

func (rbTree *Tree[K, V]) Size() int {
	return rbTree.size
}

func (rbTree *Tree[K, V]) Iter() iterator.Iter[*truple.KV[K, V]] {
	return rbTree.iter(rbTree.header.left)
}

func (rbTree *Tree[K, V]) iter(node *Node[K, V]) iterator.Iter[*truple.KV[K, V]] {
	return func() (val *truple.KV[K, V], ok bool) {
		if ok = node != nil; ok {
			val = &node.kv
			node = node.Next()
		}
		return
	}
}

func (rbTree *Tree[K, V]) iterStep(step int, node *Node[K, V]) iterator.Iter[*truple.KV[K, V]] {
	return func() (val *truple.KV[K, V], ok bool) {
		if ok = (step > 0 && node != nil); ok {
			val = &node.kv
			node = node.Next()
		}
		step--
		return
	}
}

func (rbTree *Tree[K, V]) SplitableIter() func(parallelism int) iterator.Iter[iterator.Iter[*truple.KV[K, V]]] {
	return func(parallelism int) iterator.Iter[iterator.Iter[*truple.KV[K, V]]] {
		idx := 0
		var step int
		if parallelism == 0 {
			step = rbTree.size
		} else {
			step = rbTree.size / parallelism
		}
		node := rbTree.header.left
		if parallelism <= 0 {
			return func() (iterator.Iter[*truple.KV[K, V]], bool) {
				if idx == 0 {
					idx++
					return rbTree.iter(rbTree.header.left), true
				}
				return nil, false
			}
		}
		return func() (pull iterator.Iter[*truple.KV[K, V]], ok bool) {
			if idx >= rbTree.size {
				return nil, false
			}
			i := idx
			idx += step

			n := node
			for j := i; j < i+step && node != nil; j++ {
				node = node.Next()
			}

			if i+step >= rbTree.size {
				return rbTree.iter(n), true
			}
			return rbTree.iterStep(step, n), true
		}
	}
}

func (rbTree *Tree[K, V]) Stream(parallelism ...int) stream.Stream[*truple.KV[K, V]] {
	if parallelism != nil {
		return stream.New(rbTree.SplitableIter(), parallelism[0])
	}
	return stream.New(rbTree.SplitableIter(), 0)
}

func (rbTree *Tree[K, V]) Has(key K) bool {
	if node := search(&rbTree.header, key, rbTree.comparator); node == nil {
		return false
	}
	return true
}

func (rbTree *Tree[K, V]) Get(key K) (val V) {
	node := search(&rbTree.header, key, rbTree.comparator)
	if node == nil {
		return
	}
	return node.kv.Val
}

func (rbTree *Tree[K, V]) DeleteByWholeMapping(key K, val V, valComp fn.Compare[V]) bool {
	n := search(&rbTree.header, key, rbTree.comparator)
	if n != nil && valComp(val, n.kv.Val) == 0 {
		rbTree.freeNode(rbTree.removeNode(n))
		return true
	}
	return false
}

func (rbTree *Tree[K, V]) DeleteKeys(keys ...K) int {
	changed := 0
	for i := 0; i < len(keys); i++ {
		n := rbTree.Delete(keys[i])
		if n != nil {
			rbTree.freeNode(n)
			changed++
		}
	}
	return changed
}

func (rbTree *Tree[K, V]) Delete(key K) *Node[K, V] {
	node := search(&rbTree.header, key, rbTree.comparator)
	if node == nil {
		return nil
	}
	node = remove(&rbTree.header, node, rbTree.comparator)
	// 释放tmp1节点
	// 节点数量减一
	rbTree.size--
	if rbTree.size == 0 {
		rbTree.header.right = nil
		rbTree.header.left = nil
	}
	return node
}

func (rbTree *Tree[K, V]) removeNode(node *Node[K, V]) *Node[K, V] {
	if node == nil {
		return nil
	}
	node = remove(&rbTree.header, node, rbTree.comparator)
	// 释放tmp1节点
	// 节点数量减一
	rbTree.size--
	if rbTree.size == 0 {
		rbTree.header.right = nil
		rbTree.header.left = nil
	}
	return node
}

func (rbTree *Tree[K, V]) Put(key K, val V) bool {
	if rbTree.put(key, val) {
		rbTree.size++
		return true
	}
	return false
}

func (rbTree *Tree[K, V]) PutNode(node *Node[K, V]) bool {
	var parent *Node[K, V]
	var comp int
	header := &rbTree.header
	tmpNode := header.parent
	compFunc := rbTree.comparator

	for tmpNode != nil {
		parent = tmpNode
		comp := compFunc(node.kv.Key, tmpNode.kv.Key)
		// key < tmp.kv.Key
		if comp == -1 {
			tmpNode = tmpNode.left
			comp = -1
			// key > tmp.kv.Key
		} else if comp == 1 {
			tmpNode = tmpNode.right
			comp = 1
		} else {
			tmpNode.kv.Val = node.kv.Val
			return false
		}
	}
	//return parent, comp
	// node := rbTree.allocNode(key, val) //
	// Making searched node as the header node of target node.
	node.parent = parent
	// Determining whether the searched node is the sentinel node , "node" is the header node.
	if parent == nil {
		header.left = node
		header.right = node
		header.parent = node
		//node.parent = header
		// 	The next step is determine whether "node" is the left node or right node of the parent node.
	} else if comp < 0 {
		if header.left == parent {
			header.left = node
		}
		parent.left = node
	} else {
		if header.right == parent {
			header.right = node
		}
		parent.right = node
	}

	// Making the node's left node and right node is NULL.
	node.left = nil
	node.right = nil
	// The color of new node is red
	node.color = Red
	// Entering the insert fix up stage
	putRebalanced(header, node)
	return true
}

// putNode insert a new Node into the read-black tree
func (rbTree *Tree[K, V]) put(key K, val V) bool {
	var parent *Node[K, V]
	var comp int
	header := &rbTree.header
	tmpNode := header.parent
	compFunc := rbTree.comparator

	for tmpNode != nil {
		parent = tmpNode
		comp := compFunc(key, tmpNode.kv.Key)
		// key < tmp.kv.Key
		if comp == -1 {
			tmpNode = tmpNode.left
			comp = -1
			// key > tmp.kv.Key
		} else if comp == 1 {
			tmpNode = tmpNode.right
			comp = 1
		} else {
			tmpNode.kv.Val = val
			return false
		}
	}
	//return parent, comp
	node := rbTree.allocNode(key, val) //
	// Making searched node as the header node of target node.
	node.parent = parent
	// Determining whether the searched node is the sentinel node , "node" is the header node.
	if parent == nil {
		header.left = node
		header.right = node
		header.parent = node
		//node.parent = header
		// 	The next step is determine whether "node" is the left node or right node of the parent node.
	} else if comp < 0 {
		if header.left == parent {
			header.left = node
		}
		parent.left = node
	} else {
		if header.right == parent {
			header.right = node
		}
		parent.right = node
	}

	// Making the node's left node and right node is NULL.
	node.left = nil
	node.right = nil
	// The color of new node is red
	node.color = Red
	// Entering the insert fix up stage
	putRebalanced(header, node)
	return true
}

func putRebalanced[K, V any](header, node *Node[K, V]) {
	for node != header.parent && node.parent.color == Red {
		// if node.parent == nil {
		// 	fmt.Println(node, header)
		// 	panic("AAAA")
		// }
		// if node.parent.parent == nil {
		// 	fmt.Printf("%p %p\n", node.parent, header.parent)
		// 	panic("BBBB")
		// }
		if node.parent == node.parent.parent.left {
			uncleNode := node.parent.parent.right
			if uncleNode != nil && uncleNode.color == Red {
				node.parent.color = Black
				uncleNode.color = Black
				node.parent.parent.color = Red
				node = node.parent.parent
			} else {
				if node == node.parent.right {
					node = node.parent
					leftRotate(header, node)
				}
				node.parent.color = Black
				node.parent.parent.color = Red
				rightRotate(header, node.parent.parent)
			}
		} else {
			uncleNode := node.parent.parent.left
			if uncleNode != nil && uncleNode.color == Red {
				node.parent.color = Black
				uncleNode.color = Black
				node.parent.parent.color = Red
				node = node.parent.parent

			} else {
				if node == node.parent.left {
					node = node.parent
					rightRotate(header, node)
				}
				node.parent.color = Black
				node.parent.parent.color = Red
				leftRotate(header, node.parent.parent)
			}
		}
	}
	header.parent.color = Black
}

func search[K, V any](header *Node[K, V], key K, compFunc fn.Compare[K]) *Node[K, V] {
	node := header.parent
	for node != nil {
		comp := compFunc(key, node.kv.Key)
		if comp == -1 {
			node = node.left
		} else if comp == 1 {
			node = node.right
		} else {
			return node
		}
	}
	return nil
}

func remove[K, V any](header, node *Node[K, V], compFunc fn.Compare[K]) *Node[K, V] {

	if compFunc(node.kv.Key, header.kv.Key) > -1 {
		header.right = node.Prev()
	} else if compFunc(header.left.kv.Key, node.kv.Key) > -1 {
		header.left = node.Next()
	}
	// if !compFunc(n.kv.Key, header.kv.Key) {
	// 	header.right = node.Prev()
	// } else if !compFunc(header.left.kv.Key, n.kv.Key) {
	// 	header.left = node.Next()
	// }

	var removeNode, removeNodeChild *Node[K, V] = nil, nil

	if node.left == nil || node.right == nil {
		removeNode = node
	} else {
		removeNode = node.Next()
	}
	if removeNode.left != nil {
		removeNodeChild = removeNode.left
	} else if removeNode.right != nil {
		removeNodeChild = removeNode.right
	}

	removeNodeParent := removeNode.parent
	if removeNodeChild != nil {
		removeNodeChild.parent = removeNodeParent
	}

	if removeNode.parent == nil {
		header.parent = removeNodeChild
	} else if removeNode == removeNode.parent.left {
		removeNode.parent.left = removeNodeChild
	} else {
		removeNode.parent.right = removeNodeChild
	}

	if removeNode != node {
		// 赋值
		node.kv.Key = removeNode.kv.Key
		node.kv.Val = removeNode.kv.Val
	}

	if removeNode.color == Black {
		deleteRebalanced(header, removeNodeChild, removeNodeParent)
	}

	return removeNode
}

func deleteRebalanced[K, V any](header, node, parent *Node[K, V]) {
	for node != header.parent && (node == nil || node.color == Black) {
		if node != nil {
			parent = node.parent
		}
		if node == parent.left {
			brotherNode := parent.right
			if brotherNode.color == Red {
				brotherNode.color = Black
				parent.color = Red
				leftRotate(header, parent)
				brotherNode = parent.right
			}

			if (brotherNode.left == nil || brotherNode.left.color == Black) &&
				(brotherNode.right == nil || brotherNode.right.color == Black) {
				brotherNode.color = Red
				node = parent
			} else {
				if brotherNode.right == nil || brotherNode.right.color == Black {
					if brotherNode.left != nil {
						brotherNode.left.color = Black
					}
					brotherNode.color = Red
					rightRotate(header, brotherNode)
					brotherNode = parent.right
				}
				brotherNode.color = parent.color
				parent.color = Black
				if brotherNode.right != nil {
					brotherNode.right.color = Black
				}
				leftRotate(header, parent)
				node = header.parent
			}
		} else {
			brotherNode := parent.left
			if brotherNode.color == Red {
				brotherNode.color = Black
				parent.color = Red
				rightRotate(header, parent)
				brotherNode = parent.left
			}
			if (brotherNode.left == nil || brotherNode.left.color == Black) &&
				(brotherNode.right == nil || brotherNode.right.color == Black) {
				// if getColor[T, V](brotherNode.left) == Black && getColor[T, V](brotherNode.right) == Black {
				brotherNode.color = Red
				node = parent
			} else {
				if brotherNode.left == nil || brotherNode.left.color == Black {
					// if getColor[T, V](brotherNode.left) == Black {
					if brotherNode.right != nil {
						brotherNode.right.color = Black
					}
					brotherNode.color = Red
					leftRotate(header, brotherNode)
					brotherNode = parent.left
				}
				brotherNode.color = parent.color
				parent.color = Black
				if brotherNode.left != nil {
					brotherNode.left.color = Black
				}
				rightRotate(header, parent)
				node = header.parent
			}
		}
	}
	if node != nil {
		node.color = Black
	}
}
