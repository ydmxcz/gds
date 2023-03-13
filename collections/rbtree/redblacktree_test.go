package rbtree

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/ydmxcz/gds/fn"
)

func constructTree(n int) *Tree[int, string] {
	rbTree := New[int, string](fn.Comp[int])
	for i := 0; i < n; i++ {
		rbTree.Put(i, strconv.Itoa(i))
	}
	return rbTree
}

func TestReaBlackTree_Iter(t *testing.T) {
	tree := constructTree(5000)
	iter := tree.Iter()
	count := 0
	for _, ok := iter(); ok; _, ok = iter() {
		// fmt.Print(v.Key, "  ")
		count++
	}
	fmt.Println(count)
}

func TestRbTree_SegmentedIter(t *testing.T) {
	l := constructTree(20)
	iter := l.SplitableIter()(4)
	for iter1, ok := iter(); ok; iter1, ok = iter() {
		for v1, o1 := iter1(); o1; v1, o1 = iter1() {
			fmt.Print(v1.Key, ",")
		}
		fmt.Println()
	}
}

func debugPrintTree[K any, V any](rbTree *Tree[K, V]) {
	debugPrintNode(rbTree.header.parent, 0)
}

func debugPrintNode[K any, V any](node *Node[K, V], height int) {
	if node != nil {
		for i := 0; i < height-1; i++ {
			fmt.Print("\t")
		}
		if height > 0 {
			fmt.Print("|——")
		}
		fmt.Printf("{%v:%v:", node.kv.Key, node.kv.Val)
		if node.color == Red {
			fmt.Printf("red}\n")
		} else {
			fmt.Printf("black}\n")
		}
		debugPrintNode(node.left, height+1)
		debugPrintNode(node.right, height+1)
	}
}

func TestMap(t *testing.T) {
	tree := constructTree(1000000)
	// var i collections.Map[int, string] = tree
	fmt.Println(tree.Size())
}
