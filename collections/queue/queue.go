package queue

import (
	"sync"
	"sync/atomic"
)

type Node[T any] struct {
	val  T
	next *Node[T]
}

func (n *Node[T]) Value() T {
	return n.val
}

func NewNode[T any](val T) *Node[T] {
	return &Node[T]{
		val:  val,
		next: nil,
	}
}

func New[T any]() *Linked[T] {
	return &Linked[T]{}
}

type Linked[T any] struct {
	lock  sync.Mutex
	head  Node[T]
	tail  *Node[T]
	count int64
}

func (q *Linked[T]) Push(val T) {
	q.PushNode(NewNode(val))
}

func (q *Linked[T]) Peek() T {
	return q.PeekNode().Value()
}

func (q *Linked[T]) Pop() T {
	return q.PopNode().Value()
}

func (q *Linked[T]) PushNode(node *Node[T]) {
	q.lock.Lock()
	defer q.lock.Unlock()
	if q.head.next == nil {
		q.head.next = node
		q.tail = node
	} else {
		q.tail.next = node
		q.tail = node
	}
	q.count++
}
func (q *Linked[T]) PeekNode() *Node[T] {
	q.lock.Lock()
	defer q.lock.Unlock()
	return q.head.next
}

func (q *Linked[T]) PopNode() *Node[T] {
	q.lock.Lock()
	defer q.lock.Unlock()
	n := q.head.next
	q.head.next = n.next
	if q.head.next == nil {
		q.tail = nil
	}
	q.count--
	return n
}

// Len returns the length of the queue.
func (q *Linked[T]) Len() int {
	return int(atomic.LoadInt64(&q.count))
}

// Empty returns the queue wether empty
func (q *Linked[T]) Empty() bool {
	return atomic.LoadInt64(&q.count) == 0
}
