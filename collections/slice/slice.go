package slice

import (
	"github.com/ydmxcz/gds/fn"
	"github.com/ydmxcz/gds/iterator"
	"github.com/ydmxcz/gds/stream"
	"github.com/ydmxcz/gds/util/constraints"
)

type Slice[T any] []T

func New[T any](length ...int) Slice[T] {
	if len(length) != 0 {
		return make(Slice[T], 0, length[0])
	}
	return make(Slice[T], 0, 8)
}

func Of[T any](elems ...T) Slice[T] {
	return elems
}

func (s Slice[T]) PushBack(elems ...T) {
	s = append(s, elems...)
}

func (s Slice[T]) Iter() iterator.Iter[T] {
	var idx = 0
	return func() (val T, ok bool) {
		ok = idx < len(s)
		if ok {
			val = s[idx]
			idx++
		}
		return
	}
}

func (s Slice[T]) All(yelid fn.Predicate[T]) {
	for i := 0; i < len(s); i++ {
		if !yelid(s[i]) {
			return
		}
	}
}

func (s Slice[T]) SplitableIter() func(parallelism int) iterator.Iter[iterator.Iter[T]] {
	return func(parallelism int) iterator.Iter[iterator.Iter[T]] {
		idx := 0
		var step int
		if parallelism == 0 {
			step = len(s)
		} else {
			step = len(s) / parallelism
		}
		if parallelism <= 0 {
			return func() (iterator.Iter[T], bool) {
				if idx == 0 {
					idx++
					return Slice[T](s).Iter(), true
				}
				return nil, false
			}
		}

		return func() (pull iterator.Iter[T], ok bool) {
			if idx >= len(s) {
				return nil, false
			}
			i := idx
			idx += step
			if i+step >= len(s) {
				return Slice[T](s[i:]).Iter(), true
			}
			return Slice[T](s[i : i+step]).Iter(), true
		}
	}
}

func (s Slice[T]) Stream(parallelism ...int) stream.Stream[T] {
	if parallelism != nil {
		return stream.New(s.SplitableIter(), parallelism[0])
	}
	return stream.New(s.SplitableIter(), 0)
}

// Equal reports whether two slices are equal: the same length and all
// elements equal. If the lengths are different, Equal returns false.
// Otherwise, the elements are compared in increasing index order, and the
// comparison stops at the first unequal pair.
// Floating point NaNs are not considered equal.
func Equal[T comparable](s1, s2 Slice[T]) bool {
	if len(s1) != len(s2) {
		return false
	}
	for i := range s1 {
		if s1[i] != s2[i] {
			return false
		}
	}
	return true
}

// EqualFunc reports whether two slices are equal using a comparison
// function on each pair of elements. If the lengths are different,
// EqualFunc returns false. Otherwise, the elements are compared in
// increasing index order, and the comparison stops at the first index
// for which eq returns false.
func EqualFunc[E1, E2 any](s1 Slice[E1], s2 Slice[E2], eq func(E1, E2) bool) bool {
	if len(s1) != len(s2) {
		return false
	}
	for i, v1 := range s1 {
		v2 := s2[i]
		if !eq(v1, v2) {
			return false
		}
	}
	return true
}

// Compare compares the elements of s1 and s2.
// The elements are compared sequentially, starting at index 0,
// until one element is not equal to the other.
// The result of comparing the first non-matching elements is returned.
// If both slices are equal until one of them ends, the shorter slice is
// considered less than the longer one.
// The result is 0 if s1 == s2, -1 if s1 < s2, and +1 if s1 > s2.
// Comparisons involving floating point NaNs are ignored.
func Compare[T constraints.Ordered](s1, s2 Slice[T]) int {
	s2len := len(s2)
	for i, v1 := range s1 {
		if i >= s2len {
			return +1
		}
		v2 := s2[i]
		switch {
		case v1 < v2:
			return -1
		case v1 > v2:
			return +1
		}
	}
	if len(s1) < s2len {
		return -1
	}
	return 0
}

// CompareFunc is like Compare but uses a comparison function
// on each pair of elements. The elements are compared in increasing
// index order, and the comparisons stop after the first time cmp
// returns non-zero.
// The result is the first non-zero result of cmp; if cmp always
// returns 0 the result is 0 if len(s1) == len(s2), -1 if len(s1) < len(s2),
// and +1 if len(s1) > len(s2).
func CompareFunc[E1, E2 any](s1 []E1, s2 []E2, cmp func(E1, E2) int) int {
	s2len := len(s2)
	for i, v1 := range s1 {
		if i >= s2len {
			return +1
		}
		v2 := s2[i]
		if c := cmp(v1, v2); c != 0 {
			return c
		}
	}
	if len(s1) < s2len {
		return -1
	}
	return 0
}

// Index returns the index of the first occurrence of v in s,
// or -1 if not present.
func Index[T comparable](s Slice[T], v T) int {
	for i, vs := range s {
		if v == vs {
			return i
		}
	}
	return -1
}

// Index returns the first index i satisfying f(s[i]),
// or -1 if none do.
func (s Slice[T]) Index(f func(T) bool) int {
	for i, v := range s {
		if f(v) {
			return i
		}
	}
	return -1
}

// Contains reports whether v is present in s.
func Contains[T comparable](s Slice[T], v T) bool {
	return Index(s, v) >= 0
}

// Contains reports whether at least one
// element e of s satisfies f(e).
func (s Slice[T]) Contains(f func(T) bool) bool {
	return s.Index(f) >= 0
}

// Insert inserts the values v... into s at index i,
// returning the modified slice.
// In the returned slice r, r[i] == v[0].
// Insert panics if i is out of range.
// This function is O(len(s) + len(v)).
func Insert[S ~[]E, E any](s S, i int, v ...E) S {
	tot := len(s) + len(v)
	if tot <= cap(s) {
		s2 := s[:tot]
		copy(s2[i+len(v):], s[i:])
		copy(s2[i:], v)
		return s2
	}
	s2 := make(S, tot)
	copy(s2, s[:i])
	copy(s2[i:], v)
	copy(s2[i+len(v):], s[i:])
	return s2
}
func (s Slice[T]) Insert(i int, v ...T) {
	s = Insert(s, i, v...)
}

// Delete removes the elements s[i:j] from s, returning the modified slice.
// Delete panics if s[i:j] is not a valid slice of s.
// Delete modifies the contents of the slice s; it does not create a new slice.
// Delete is O(len(s)-j), so if many items must be deleted, it is better to
// make a single call deleting them all together than to delete one at a time.
// Delete might not modify the elements s[len(s)-(j-i):len(s)]. If those
// elements contain pointers you might consider zeroing those elements so that
// objects they reference can be garbage collected.
func Delete[S ~[]E, E any](s S, i, j int) S {
	_ = s[i:j] // bounds check

	return append(s[:i], s[j:]...)
}

// Replace replaces the elements s[i:j] by the given v, and returns the
// modified slice. Replace panics if s[i:j] is not a valid slice of s.
func Replace[S ~[]E, E any](s S, i, j int, v ...E) S {
	_ = s[i:j] // verify that i:j is a valid subslice
	tot := len(s[:i]) + len(v) + len(s[j:])
	if tot <= cap(s) {
		s2 := s[:tot]
		copy(s2[i+len(v):], s[j:])
		copy(s2[i:], v)
		return s2
	}
	s2 := make(S, tot)
	copy(s2, s[:i])
	copy(s2[i:], v)
	copy(s2[i+len(v):], s[j:])
	return s2
}

// Clone returns a copy of the slice.
// The elements are copied using assignment, so this is a shallow clone.
func Clone[S ~[]E, E any](s S) S {
	// Preserve nil in case it matters.
	if s == nil {
		return nil
	}
	return append(S([]E{}), s...)
}

// Compact replaces consecutive runs of equal elements with a single copy.
// This is like the uniq command found on Unix.
// Compact modifies the contents of the slice s; it does not create a new slice.
// When Compact discards m elements in total, it might not modify the elements
// s[len(s)-m:len(s)]. If those elements contain pointers you might consider
// zeroing those elements so that objects they reference can be garbage collected.
func Compact[S ~[]E, E comparable](s S) S {
	if len(s) < 2 {
		return s
	}
	i := 1
	last := s[0]
	for _, v := range s[1:] {
		if v != last {
			s[i] = v
			i++
			last = v
		}
	}
	return s[:i]
}

// CompactFunc is like Compact but uses a comparison function.
func CompactFunc[S ~[]E, E any](s S, eq func(E, E) bool) S {
	if len(s) < 2 {
		return s
	}
	i := 1
	last := s[0]
	for _, v := range s[1:] {
		if !eq(v, last) {
			s[i] = v
			i++
			last = v
		}
	}
	return s[:i]
}

// Grow increases the slice's capacity, if necessary, to guarantee space for
// another n elements. After Grow(n), at least n elements can be appended
// to the slice without another allocation. If n is negative or too large to
// allocate the memory, Grow panics.
func Grow[S ~[]E, E any](s S, n int) S {
	if n < 0 {
		panic("cannot be negative")
	}
	if n -= cap(s) - len(s); n > 0 {
		// TODO(https://go.dev/issue/53888): Make using []E instead of S
		// to workaround a compiler bug where the runtime.growslice optimization
		// does not take effect. Revert when the compiler is fixed.
		s = append([]E(s)[:cap(s)], make([]E, n)...)[:len(s)]
	}
	return s
}

// Clip removes unused capacity from the slice, returning s[:len(s):len(s)].
func Clip[S ~[]E, E any](s S) S {
	return s[:len(s):len(s)]
}
