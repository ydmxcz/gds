package priorityqueue

import (
	"github.com/ydmxcz/gds/fn"
	"github.com/ydmxcz/gds/util/constraints"
)

type Queue[T any] struct {
	elements []T
	comp     fn.Compare[T]
	size     int
}

func New[T any](lessFunc fn.Compare[T], size ...int) *Queue[T] {
	q := &Queue[T]{}
	q.Init(lessFunc, size...)
	return q
}

func (q *Queue[T]) Init(lessFunc fn.Compare[T], size ...int) {
	s := 0
	if len(size) == 0 {
		s = size[0]
	} else {
		s = 8
	}
	q.elements = make([]T, s)
	q.comp = lessFunc
	q.size = 0
}

func OfWith[T constraints.Ordered](comp fn.Compare[T], elements ...T) *Queue[T] {
	q := New(comp)
	for _, element := range elements {
		q.Push(element)
	}
	return q
}

func (q *Queue[T]) Push(val T) bool {
	if q.size >= len(q.elements) {
		q.elements = append(q.elements, val)
	} else {
		q.elements[q.size] = val
	}
	q.size++
	q.percolateUp(q.size - 1)
	return true
}

func (q *Queue[T]) Pop() (val T, ok bool) {
	if q.size == 0 {
		return
	}
	n := q.size - 1
	q.elements[0], q.elements[n] = q.elements[n], q.elements[0]
	//val = q.deleteMin()
	q.percolateDown(0, n)
	q.size--
	val = q.elements[q.size]
	return val, true
}

func (q *Queue[T]) percolateUp(pos int) {
	for {
		i := (pos - 1) / 2 // parent
		if i == pos || q.comp(q.elements[pos], q.elements[i]) >= 0 {
			break
		}
		q.elements[i], q.elements[pos] = q.elements[pos], q.elements[i]
		pos = i
	}
}

func (q *Queue[T]) percolateDown(i0, n int) bool {
	i := i0
	for {
		j1 := (i << 1) + 1
		if j1 >= n || j1 < 0 { // j1 < 0 after int overflow
			break
		}
		j := j1 // left child
		if j2 := j1 + 1; j2 < n && q.comp(q.elements[j2], q.elements[j1]) < 0 {
			j = j2 // = 2*i + 2  // right child
		}
		if q.comp(q.elements[j], q.elements[i]) >= 0 {
			break
		}
		//h.Swap(i, j)
		q.elements[i], q.elements[j] = q.elements[j], q.elements[i]
		i = j
	}
	return i > i0
}

func (q *Queue[T]) Len() int {
	return q.size
}

func (q *Queue[T]) Peek() (val T, ok bool) {
	if q.size > 0 {
		return q.elements[0], true
	}
	return
}
