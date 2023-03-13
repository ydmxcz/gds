package bplustree

import (
	"fmt"

	"github.com/ydmxcz/gds/collections/truple"
)

type Node[K, V any] struct {
	//parent   *Node[K, V]
	prev     *Node[K, V]
	next     *Node[K, V]
	values   []truple.KV[K, V]
	children []*Node[K, V]
	num      int
	isLeaf   bool
}

func Inorder[K, V any](node *Node[K, V]) {
	if nil != node {
		Inorder(node.children[0])
		for i := 0; i < node.num; i++ {
			fmt.Println(node.values[i])
			Inorder(node.children[i+1])
		}
	}
}

func (node *Node[K, V]) MinimumNode() *Node[K, V] {
	y := node
	if y == nil {
		return nil
	}
	for !y.isLeaf {
		y = y.children[0]
	}
	return y
}

func (node *Node[K, V]) MaximumNode() *Node[K, V] {
	y := node
	for !y.isLeaf {
		y = y.children[y.num]
	}
	return y
}
