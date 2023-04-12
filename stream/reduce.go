package stream

func Reduce[T any](stm Stream[T], reduce func(a, b T) T) (res T) {
	res, ok := First(stm)
	if !ok {
		return
	}
	return FoldWith(stm, res, reduce, reduce)
}

func ReduceWith[T, R any](stm Stream[T], identity R, reduce func(R, T) R, combiner func(R, R) R) (res R) {

	return FoldWith(stm, identity, reduce, combiner)

}
