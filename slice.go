package vino

// SliceIter iterates over the slice s and calls the provided function f
// for each element. The function f receives the index and the corresponding
// element as arguments.
func SliceIter[T any](s []T, f func(int, T)) {
	for i := range s {
		f(i, s[i])
	}
}

// SliceWalk iterates over the slice s and calls the provided function f
// for each element. The function f receives the index and element as
// arguments and returns a boolean value. If f returns false, the iteration
// is stopped immediately.
func SliceWalk[T any](s []T, f func(int, T) bool) {
	for i := range s {
		if !f(i, s[i]) {
			return
		}
	}
}

// SliceUnique returns a new slice that contains only the unique elements
// from the input slice s. The order of elements in the resulting slice
// is not guaranteed.
func SliceUnique[T comparable](s []T) []T {
	mp := make(map[T]struct{}, len(s))
	for _, x := range s {
		mp[x] = struct{}{}
	}

	ret := make([]T, 0, len(mp))
	for x := range mp {
		ret = append(ret, x)
	}

	return ret
}

// SliceToStream converts the slice s into a Stream of type T. If a
// repeat count is provided, the resulting stream will repeat the
// slice that many times; otherwise, it defaults to no repetition.
func SliceToStream[T any](s []T, repeat ...int) Stream[T] {
	if len(repeat) > 0 {
		return NewRepeatedStream(s, repeat[0])
	}
	return NewRepeatedStream(s, 0)
}
