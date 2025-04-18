package vino

import "sort"

type Bisect[T any] interface {
	BisectLeft(xs []T, x T) int
	BisectRight(xs []T, x T) int
}

type Bisector[T any] struct {
	cmp func(T, T) int
}

func (b Bisector[T]) BisectLeft(xs []T, v T) int {
	return b.bisectLeftRange(xs, v, 0, len(xs))
}

func (b Bisector[T]) bisectLeftRange(xs []T, v T, lo, hi int) int {
	s := xs[lo:hi]
	return sort.Search(len(xs), func(i int) bool {
		return b.cmp(s[i], v) >= 0
	})
}

func (b Bisector[T]) BisectRight(xs []T, v T) int {
	return b.bisectRightRange(xs, v, 0, len(xs))
}

func (b Bisector[T]) bisectRightRange(xs []T, v T, lo, hi int) int {
	s := xs[lo:hi]
	return sort.Search(len(s), func(i int) bool {
		return b.cmp(s[i], v) > 0
	})
}

func NewBisector[T any](cmp func(T, T) int) Bisect[T] {
	return Bisector[T]{cmp}
}
