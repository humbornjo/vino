package vino

func SliceIter[T any](s []T, f func(int, T)) {
	for i := range s {
		f(i, s[i])
	}
}

func SliceWalk[T any](s []T, f func(int, T) bool) {
	for i := range s {
		if !f(i, s[i]) {
			return
		}
	}
}

func SliceHas[T comparable](s []T, x T) bool {
	for i := range s {
		if s[i] == x {
			return true
		}
	}
	return false
}

func SliceToStream[T any](s []T, magnitude ...int) Stream[T] {
	repeat := 0
	if len(magnitude) > 0 {
		repeat = magnitude[0]
	}

	return NewRepeatedStream(s, repeat)
}
