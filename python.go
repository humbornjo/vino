package vino

import (
	"sort"
)

type IfaceOmni[I any, V any] interface {
	IfaceScalp[I]
	IfaceMerge[V]
	IfaceClamp[I, V]
	IfaceBisect[I, V]
}

type OmniImpl[I any, V any] struct {
	IfaceScalp[I]
	IfaceMerge[V]
	IfaceClamp[I, V]
	IfaceBisect[I, V]
}

func NewOmniImpl[I any, V any](
	idx func(V) I,
	add func(I, int) I,
	cmp func(I, I) int,
) IfaceOmni[I, V] {
	return OmniImpl[I, V]{
		IfaceScalp:  NewScalpImpl(cmp, add),
		IfaceMerge:  NewMergeImpl(idx, cmp),
		IfaceClamp:  NewClampImpl(idx, cmp),
		IfaceBisect: NewBisectImpl(idx, cmp),
	}
}

type IfaceBisect[I any, V any] interface {
	BisectLeft(xs []V, x I) int
	BisectRight(xs []V, x I) int
}

type BisectImpl[I any, V any] struct {
	idx func(V) I
	cmp func(I, I) int
}

func (b BisectImpl[I, V]) BisectLeft(xs []V, idx I) int {
	return sort.Search(len(xs), func(i int) bool {
		return b.cmp(b.idx(xs[i]), idx) >= 0
	})
}

func (b BisectImpl[I, V]) BisectRight(xs []V, idx I) int {
	return sort.Search(len(xs), func(i int) bool {
		return b.cmp(b.idx(xs[i]), idx) > 0
	})
}

func NewBisectImpl[I any, V any](idx func(V) I, cmp func(I, I) int) IfaceBisect[I, V] {
	return BisectImpl[I, V]{idx, cmp}
}

type (
	Interval[T any] struct {
		Left  T
		Right T
	}
	IfaceScalp[T any] interface {
		Scalp(xs []Interval[T], lo T, hi T) []Interval[T]
	}
)

type ScalpImpl[T any] struct {
	cmp func(T, T) int
	add func(T, int) T
}

func (s ScalpImpl[T]) Scalp(xs []Interval[T], lo T, hi T) []Interval[T] {
	ret := make([]Interval[T], 0, len(xs))
	for _, interval := range xs {
		cmpRL := s.cmp(hi, interval.Left)
		cmpLR := s.cmp(lo, interval.Right)
		if cmpRL < 0 || cmpLR > 0 {
			ret = append(ret, interval)
			continue
		}

		cmpLL := s.cmp(lo, interval.Left)
		cmpRR := s.cmp(hi, interval.Right)
		if cmpLL > 0 && cmpRR < 0 {
			ret = append(ret, Interval[T]{Left: interval.Left, Right: s.add(lo, -1)})
			ret = append(ret, Interval[T]{Left: s.add(hi, 1), Right: interval.Right})
		} else if cmpLL > 0 && cmpLR <= 0 {
			ret = append(ret, Interval[T]{Left: interval.Left, Right: s.add(lo, -1)})
		} else if cmpRR < 0 && cmpRL >= 0 {
			ret = append(ret, Interval[T]{Left: s.add(hi, 1), Right: interval.Right})
		}
	}
	return ret
}

func NewScalpImpl[T any](cmp func(T, T) int, add func(T, int) T) IfaceScalp[T] {
	return ScalpImpl[T]{cmp, add}
}

type IfaceMerge[T any] interface {
	Merge(xs0 []T, xs1 []T) []T
}

type MergeImpl[I any, V any] struct {
	idx func(V) I
	cmp func(I, I) int
}

func (m MergeImpl[I, V]) Merge(xs0 []V, xs1 []V) []V {
	if len(xs0) == 0 {
		return xs1
	}
	if len(xs1) == 0 {
		return xs0
	}

	i, j := 0, 0
	ret := make([]V, 0, len(xs0)+len(xs1))
	for i < len(xs0) && j < len(xs1) {
		cmp := m.cmp(m.idx(xs0[i]), m.idx(xs1[j]))
		if cmp < 0 {
			ret = append(ret, xs0[i])
			i++
		} else if cmp > 0 {
			ret = append(ret, xs1[j])
			j++
		} else {
			ret = append(ret, xs0[i])
			i++
			j++
		}
	}
	for ; i < len(xs0); i++ {
		ret = append(ret, xs0[i])
	}
	for ; j < len(xs1); j++ {
		ret = append(ret, xs1[j])
	}
	return ret
}

func NewMergeImpl[I any, V any](idx func(V) I, cmp func(I, I) int) IfaceMerge[V] {
	return MergeImpl[I, V]{idx, cmp}
}

type IfaceClamp[I any, V any] interface {
	Clamp(xs []V, lo I, hi I) []V
}

type ClampImpl[I any, V any] struct {
	idx func(V) I
	cmp func(I, I) int
}

func (c ClampImpl[I, V]) Clamp(xs []V, lo I, hi I) []V {
	ret := make([]V, 0, len(xs))
	for _, x := range xs {
		if c.cmp(lo, c.idx(x)) > 0 {
			continue
		}

		if c.cmp(hi, c.idx(x)) < 0 {
			break
		}

		ret = append(ret, x)
	}
	return ret
}

func NewClampImpl[I any, V any](idx func(V) I, cmp func(I, I) int) IfaceClamp[I, V] {
	return ClampImpl[I, V]{idx, cmp}
}
