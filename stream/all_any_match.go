package stream

func AllMatch[T ~int](stm Stream[T], callback func(T) bool) bool {
	max, ok := First(stm)
	if !ok {
		return false
	}
	_, ok = TryFold(stm, max, func(a, b T) (T, bool) {
		if callback(b) {
			return a, true
		}
		return a, false
	})
	return ok
}

func AnyMatch[T ~int](stm Stream[T], callback func(T) bool) bool {
	max, ok := First(stm)
	if !ok {
		return false
	}
	_, ok = TryFold(stm, max, func(a, b T) (T, bool) {
		if callback(b) {
			return a, false
		}
		return a, true
	})
	return !ok
}
