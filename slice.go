package vino

func SliceIter[T any](s []T, f func(int, T)) {
	for i, v := range s {
		f(i, v)
	}
}

func SliceWalk[T any](s []T, f func(int, T) bool) {
	for i, v := range s {
		if !f(i, v) {
			return
		}
	}
}

func SliceHas[T comparable](s []T, x T) bool {
	for _, v := range s {
		if v == x {
			return true
		}
	}
	return false
}
