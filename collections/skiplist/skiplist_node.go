package skiplist

type Node[K, V any] struct {
	value    V
	key      K
	backward *Node[K, V]
	level    []Level[K, V]
}

type Level[K, V any] struct {
	forward *Node[K, V]
	span    int
}

func (n *Node[K, V]) IsValid() bool {
	return n != nil
}

func (n *Node[K, V]) Next() *Node[K, V] {
	return n.level[0].forward

}

func (n *Node[K, V]) GetKey() K {
	return n.key
}

func (n *Node[K, V]) GetValue() V {
	return n.value
}

func (n *Node[K, V]) SetValue(val V) {
	n.value = val
}

func (n *Node[K, V]) Prev() *Node[K, V] {
	return n.backward
}
