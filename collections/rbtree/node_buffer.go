package rbtree

type nodeBuffer[K, V any] struct {
	head           Node[K, V]
	header         *Node[K, V]
	nodes          []Node[K, V]
	sliceCursor    int
	linkedListSize int
	bufferedLength int
}

func newNodeBuffer[K, V any](size int) *nodeBuffer[K, V] {
	return newNodeBufferWith(make([]Node[K, V], size, size))
}

func newNodeBufferWith[K, V any](nodes []Node[K, V]) *nodeBuffer[K, V] {
	nb := &nodeBuffer[K, V]{}
	nb.InitWith(nodes)
	return nb
}

func (nb *nodeBuffer[K, V]) Init(size int) {
	nb.header = &nb.head
	nb.header.right = nil
	nb.linkedListSize = 0
	nb.FillWith(make([]Node[K, V], size, size), true)
}

func (nb *nodeBuffer[K, V]) InitWith(nodes []Node[K, V]) {
	nb.header = &nb.head
	nb.header.right = nil
	nb.linkedListSize = 0
	nb.FillWith(nodes, true)
}

func (nb *nodeBuffer[K, V]) IsEmpty() bool {
	return nb.Size() == 0
}

func (nb *nodeBuffer[K, V]) Size() int {
	return nb.linkedListSize + nb.sliceCursor + 1
}

func (nb *nodeBuffer[K, V]) Clear() {
	nb.header = &nb.head
	nb.header.right = nil
	nb.linkedListSize = 0
	nb.sliceCursor = -1
	nb.nodes = nil
}

func (nb *nodeBuffer[K, V]) FillN(n int) {
	if nb.sliceCursor == -1 {
		nb.FillWith(make([]Node[K, V], n), true)
	}
}

func (nb *nodeBuffer[K, V]) Fill() {
	if nb.sliceCursor == -1 {
		nb.FillWith(make([]Node[K, V], nb.bufferedLength), false)
	}
}

func (nb *nodeBuffer[K, V]) FillWith(nodes []Node[K, V], all bool) {
	if all {
		nb.bufferedLength = len(nodes)
	}
	nb.sliceCursor = nb.bufferedLength - 1
	nb.nodes = nodes[:nb.bufferedLength]
}

func (nb *nodeBuffer[K, V]) Free(node *Node[K, V]) {
	if nb.linkedListSize == nb.bufferedLength {
		return
	}
	node.right = nb.header.right
	nb.header.right = node
	nb.linkedListSize++
}

func (nb *nodeBuffer[K, V]) Alloc() (n *Node[K, V]) {

allocInSlice:
	if nb.sliceCursor != -1 {
		n = &nb.nodes[nb.sliceCursor]
		nb.sliceCursor--
		return n
	}

	// linked list is nil,filling the slice and try to reallocate
	if nb.linkedListSize == 0 {
		nb.Fill()
		goto allocInSlice
	}
	// slice is nil,but linked list is not nil,try to allocate node from head
	n = nb.header.right
	nb.header.right = n.right
	nb.linkedListSize--
	return n
}
