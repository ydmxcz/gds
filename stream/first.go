package stream

func First[T any](stm Stream[T]) (T, bool) {
	var init T
	pull, b := stm.Activate(stm.parallelism)()
	if !b {
		return init, false
	}
	// accum := init
	for val, ok := pull(); ok; val, ok = pull() {
		return val, true
	}
	return init, false
}
