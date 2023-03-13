package rbtree

import "github.com/ydmxcz/gds/collections/truple"

const (
	Red   = false
	Black = true
)

type TreeColorType bool

type Node[K any, V any] struct {
	color  TreeColorType
	parent *Node[K, V]
	left   *Node[K, V]
	right  *Node[K, V]
	kv     truple.KV[K, V]
}

func newRbTreeNodeByVal[K any, V any](key K, val V) *Node[K, V] {
	return &Node[K, V]{
		kv: truple.KV[K, V]{
			Val: val,
			Key: key,
		},
		color: Red,
	}
}

func (node *Node[K, V]) GetKey() K {
	return node.kv.Key
}

func (node *Node[K, V]) GetValue() V {
	return node.kv.Val
}

func (node *Node[K, V]) SetKey(key K) {
	node.kv.Key = key
}

func (node *Node[K, V]) SetValue(val V) {
	node.kv.Val = val
}

func (node *Node[K, V]) GetLeft() *Node[K, V] {
	return node.left
}

func (node *Node[K, V]) GetRight() *Node[K, V] {
	return node.right
}

func (node *Node[K, V]) GetParent() *Node[K, V] {
	return node.parent
}

// MinChildNode find min node
func (node *Node[K, V]) MinChildNode() *Node[K, V] {
	newNode := node
	for newNode.left != nil {
		newNode = newNode.left
	}
	return newNode
}

// MaxChildNode find max node
func (node *Node[K, V]) MaxChildNode() *Node[K, V] {
	newNode := node
	for newNode.right != nil {
		newNode = newNode.right
	}
	return newNode
}

func (node *Node[K, V]) IsValid() bool {
	return node != nil
}

func (node *Node[K, V]) Next() *Node[K, V] {
	n := node
	// The min node of `node`'s right subtree is obtained first
	if n.right != nil {
		return n.right.MinChildNode()
	}
	var y = n.parent
	for y != nil && n == y.right {
		n = y
		y = y.parent
	}
	return y
}

func (node *Node[K, V]) Prev() *Node[K, V] {
	n := node
	if n.left != nil {
		return n.left.MaxChildNode()
	} else {
		var y = n.parent
		for y != nil && y.left == n {
			n = y
			y = y.parent
		}
		return y
	}
}

// pushToNext using the `right` and `left` filed as the `next` and `previous` pointer of linked list respectively,
// the `parent` filed represents the linked list's head node which it belonged to
func (node *Node[K, V]) addToNext(n *Node[K, V]) {
	n.right = node.right
	node.right.left = n

	n.left = node
	node.right = n

}

func (node *Node[K, V]) removeSelf() *Node[K, V] {
	// if only one node in the linked list
	// this way just circular set its own filed
	node.right.left = node.left
	node.left.right = node.right
	// clear its own pointer

	return node
}

const (
	PreOrder = iota
	Inorder
	PostOrder
)

// func preorderTraversalRecursion[K any, V any](node *Node[K, V],
// 	callBack function.Consumer[*Node[K, V]]) {
// 	if node != nil {
// 		callBack(node)
// 		preorderTraversalRecursion(node.left, callBack)
// 		preorderTraversalRecursion(node.right, callBack)
// 	}
// }

// func inorderTraversalRecursion[K any, V any](node *Node[K, V],
// 	callBack function.Consumer[*Node[K, V]]) {
// 	if node != nil {
// 		inorderTraversalRecursion(node.left, callBack)
// 		callBack(node)
// 		inorderTraversalRecursion(node.right, callBack)
// 	}
// }

// func postorderTraversalRecursion[K any, V any](node *Node[K, V],
// 	callBack function.Consumer[*Node[K, V]]) {
// 	if node != nil {
// 		postorderTraversalRecursion(node.left, callBack)
// 		postorderTraversalRecursion(node.right, callBack)
// 		callBack(node)
// 	}
// }
