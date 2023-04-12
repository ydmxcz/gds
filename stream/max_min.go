package stream

func Max[T ~int](stm Stream[T]) (max T) {
	return MaxBy(stm, func(a, b T) T {
		if a > b {
			return a
		}
		return b
	})
}

func MaxBy[T ~int](stm Stream[T], comp func(a, b T) T) (max T) {
	max, ok := First(stm)
	if !ok {
		return
	}
	return FoldWith(stm, max, comp, comp)
}

func Min[T ~int](stm Stream[T]) (min T) {
	return MinBy(stm, func(a, b T) T {
		if a < b {
			return a
		}
		return b
	})
}

func MinBy[T ~int](stm Stream[T], comp func(a, b T) T) (max T) {
	max, ok := First(stm)
	if !ok {
		return
	}
	return FoldWith(stm, max, comp, comp)

}
