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

type sliceStream[T any] struct {
	s []T
	i int
}

func (s sliceStream[T]) Next() (T, bool) {
	idx := s.i
	if idx < len(s.s) {
		v := s.s[idx]
		s.i = idx + 1
		return v, true
	}
	return *new(T), false
}

func SliceAsStream[T any](s []T) Stream[T] {
	return sliceStream[T]{s, 0}
}
