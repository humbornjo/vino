package vino

func SliceIter[T any](s []T, f func(T) bool) {
	for _, v := range s {
		if !f(v) {
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
