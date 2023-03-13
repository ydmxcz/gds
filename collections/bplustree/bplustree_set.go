package bplustree

// type TreeSet[T any] struct {
// 	tree Tree[T, constrains.Void]
// }

// func NewSet[T any](compareFunc comparison.OneSided[T]) *TreeSet[T] {
// 	t := &TreeSet[T]{}
// 	t.tree.Init(10, compareFunc)
// 	return t
// }

// func Of[T constrains.Ordered](elems ...T) collections.Set[T] {
// 	set := NewSet(comparison.Less[T])
// 	for i := 0; i < len(elems); i++ {
// 		set.Add(elems[i])
// 	}
// 	return set
// }

// func OfWith[T any](compFunc comparison.OneSided[T], elements ...T) *TreeSet[T] {
// 	newSet := NewSet(compFunc)
// 	for _, element := range elements {
// 		newSet.Add(element)
// 	}
// 	return newSet
// }

// func Intersection[T any](set1, set2 *TreeSet[T]) *TreeSet[T] {
// 	newSet := NewSet(set1.tree.GetCompareFunc())
// 	//algorithm.IntersectionBetweenTwoSet[T](set1.Iter(), set2, newSet)
// 	return newSet
// }

// func Union[T any](set1, set2 *TreeSet[T]) *TreeSet[T] {
// 	newSet := NewSet(set1.tree.GetCompareFunc())
// 	//algorithm.UnionBetweenTwoSet[T](set1.Iter(), set2.Iter(), newSet)
// 	return newSet
// }

// func Different[T any](set1, set2 *TreeSet[T]) *TreeSet[T] {
// 	newSet := NewSet(set1.tree.GetCompareFunc())
// 	//algorithm.DifferentBetweenTwoSet[T](set1, set1.Iter(), set2, set2.Iter(), newSet)
// 	return newSet
// }

// func (ts *TreeSet[T]) AddIter(base iterator.Base[T]) bool {
// 	return ts.Add(base.GetValue())
// }

// func (ts *TreeSet[T]) AddAllIter(first, end iterator.Forward[T]) int {
// 	changed := 0
// 	for ; first.Comp(end); first.Next() {
// 		if ts.Add(first.GetValue()) {
// 			changed++
// 		}
// 	}
// 	return changed
// }

// func (ts *TreeSet[T]) DeleteAllIter(first, end iterator.Forward[T]) int {
// 	changed := 0
// 	for ; first.Comp(end); first.Next() {
// 		if ts.Delete(first.GetValue()) {
// 			changed++
// 		}
// 	}
// 	return changed
// }

// func (ts *TreeSet[T]) DeleteIter(base iterator.Base[T]) bool {
// 	return ts.Delete(base.GetValue())
// }

// func (ts *TreeSet[T]) String() string {
// 	var buf bytes.Buffer

// 	buf.WriteString("[")
// 	n := ts.Iter()
// 	for n.IsValid() {
// 		buf.WriteString(fmt.Sprintf("%v", n.GetValue()))
// 		n.Next()
// 		if n.IsValid() {
// 			buf.WriteByte(' ')
// 		} else {
// 			break
// 		}
// 	}
// 	buf.WriteByte(']')
// 	return buf.String()
// }

// func (ts *TreeSet[T]) Size() int {
// 	return ts.tree.Size()
// }

// func (ts *TreeSet[T]) IsEmpty() bool {
// 	return ts.tree.IsEmpty()
// }

// func (ts *TreeSet[T]) Contains(elem T) bool {
// 	return ts.tree.GetEntry(elem) != nil
// }

// func (ts *TreeSet[T]) End() iterator.Forward[T] {
// 	iter := InitUnaryIterator(&ts.tree)
// 	return &iter
// }

// func (ts *TreeSet[T]) Iter() iterator.Forward[T] {
// 	iter := InitUnaryIterator(&ts.tree)
// 	return &iter

// }

// func (ts *TreeSet[T]) RIter() iterator.Forward[T] {
// 	iter := InitReverseUnaryIterator(&ts.tree)
// 	return &iter
// }

// func (ts *TreeSet[T]) BidIter() iterator.Bidirectional[T] {
// 	iter := InitUnaryIterator(&ts.tree)
// 	return &iter
// }

// func (ts *TreeSet[T]) ToSlice() []T {
// 	size := ts.tree.Size()
// 	slice := make([]T, size)
// 	elem := ts.tree.Iter()
// 	for i := 0; i < size && elem.IsValid(); i++ {
// 		slice[i] = elem.GetValue().GetKey()
// 		elem.Next()
// 	}
// 	return slice
// }

// func (ts *TreeSet[T]) Add(elem T) bool {
// 	return ts.tree.Put(elem, constrains.Void{})
// }

// func (ts *TreeSet[T]) AddAll(elems ...T) int {
// 	changed := 0
// 	for _, elem := range elems {
// 		if ts.tree.Put(elem, constrains.Void{}) {
// 			changed++
// 		}
// 	}
// 	return changed
// }

// func (ts *TreeSet[T]) Delete(elem T) bool {
// 	return ts.tree.Delete(elem)
// }

// func (ts *TreeSet[T]) DeleteAll(elems ...T) int {
// 	changed := 0
// 	for i := 0; i < len(elems); i++ {
// 		if ts.tree.Delete(elems[i]) {
// 			changed++
// 		}
// 	}
// 	return changed
// }

// func (ts *TreeSet[T]) Clear() {
// 	ts.tree.Clear()
// }
