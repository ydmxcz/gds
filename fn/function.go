package fn

type Function[T, R any] func(T) R

func Compose[T, V, R any](before Function[V, T], curr Function[T, R]) Function[V, R] {
	return func(v V) R {
		return curr(before(v))
	}
}

func AndThen[T, V, R any](curr Function[T, R], after Function[R, V]) Function[T, V] {
	return func(t T) V {
		return after(curr(t))
	}
}

type Predicate[T any] func(T) bool

func (p Predicate[T]) And(other Predicate[T]) Predicate[T] {
	return func(t T) bool {
		return p(t) && other(t)
	}
}

func (p Predicate[T]) Negate() Predicate[T] {
	return func(t T) bool {
		return !p(t)
	}
}

func (p Predicate[T]) Or(other Predicate[T]) Predicate[T] {
	return func(t T) bool {
		return p(t) || other(t)
	}
}

func Not[T any](target Predicate[T]) Predicate[T] {
	return target.Negate()
}

type Supplier[T any] func() T

type Consumer[T any] func(T)

func (c Consumer[T]) AndThen(after Consumer[T]) Consumer[T] {
	return func(t T) {
		c(t)
		after(t)
	}
}

// Represents an operation that accepts two input arguments and returns no result.
// This is the two-arity specialization of Consumer.
// Unlike most other functional interfaces,
// BiConsumer is expected to operate via side-effects.
type BinConsumer[T, U any] func(T, U)

func (curr BinConsumer[T, U]) AndThen(after BinConsumer[T, U]) BinConsumer[T, U] {
	return func(t T, u U) {
		curr(t, u)
		after(t, u)
	}
}

// Represents a function that accepts two arguments and produces a result.
// This is the two-arity specialization of Function.
type BinFunction[T, U, R any] func(T, U) R

type BinOperator[T any] BinFunction[T, T, T]

type BinPredicate[T, U any] func(T, U) bool

func (curr BinPredicate[T, U]) And(other BinPredicate[T, U]) BinPredicate[T, U] {
	return func(t T, u U) bool {
		return curr(t, u) && other(t, u)
	}
}

func (curr BinPredicate[T, U]) Or(other BinPredicate[T, U]) BinPredicate[T, U] {
	return func(t T, u U) bool {
		return curr(t, u) && other(t, u)
	}
}

func (curr BinPredicate[T, U]) negate() BinPredicate[T, U] {
	return func(t T, u U) bool {
		return !curr(t, u)
	}
}

type Pull[T any] func() (T, bool)

type Push[T any] func(Predicate[T])

type PushPred[T any] Predicate[Predicate[T]]
