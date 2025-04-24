package vino_test

import (
	"testing"

	. "github.com/humbornjo/vino"
	"github.com/stretchr/testify/assert"
)

func cmpInt(a, b int) int {
	return a - b
}

func addInt(a, d int) int {
	return a + d
}

func idxInt(v int) int {
	return v
}

func TestBisectorInt(t *testing.T) {
	b := NewBisectImpl(idxInt, cmpInt)

	tests := []struct {
		name      string
		xs        []int
		x         int
		wantLeft  int
		wantRight int
	}{
		{
			"empty slice",
			[]int{}, 5, 0, 0,
		},
		{
			"single less",
			[]int{3}, 5, 1, 1,
		},
		{
			"single equal",
			[]int{5}, 5, 0, 1,
		},
		{
			"single greater",
			[]int{7}, 5, 0, 0,
		},
		{
			"multiple no dup",
			[]int{1, 3, 5, 7}, 5, 2, 3,
		},
		{
			"multiple left end",
			[]int{1, 3, 5, 7}, 0, 0, 0,
		},
		{
			"multiple right end",
			[]int{1, 3, 5, 7}, 8, 4, 4,
		},
		{
			"multiple between",
			[]int{1, 3, 5, 7}, 4, 2, 2,
		},
		{
			"duplicates",
			[]int{1, 3, 3, 3, 5}, 3, 1, 4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			left := b.BisectLeft(tt.xs, tt.x)
			right := b.BisectRight(tt.xs, tt.x)
			assert.Equal(t, tt.wantLeft, left, "BisectLeft(%v, %d)", tt.xs, tt.x)
			assert.Equal(t, tt.wantRight, right, "BisectRight(%v, %d)", tt.xs, tt.x)
		})
	}
}

func TestScalperInt(t *testing.T) {
	s := NewScalpImpl(cmpInt, addInt)

	tests := []struct {
		name   string
		xs     []Interval[int]
		lo, hi int
		want   []Interval[int]
	}{
		{
			"no overlap below",
			[]Interval[int]{{1, 3}},
			5, 7,
			[]Interval[int]{{1, 3}},
		},
		{
			"no overlap above",
			[]Interval[int]{{10, 15}},
			1, 5,
			[]Interval[int]{{10, 15}},
		},
		{
			"full contained interval removed",
			[]Interval[int]{{5, 10}},
			1, 15,
			[]Interval[int]{},
		},
		{
			"partial left cut",
			[]Interval[int]{{1, 10}},
			5, 8,
			[]Interval[int]{{1, 4}, {9, 10}},
		},
		{
			"partial right cut",
			[]Interval[int]{{1, 10}},
			3, 6,
			[]Interval[int]{{1, 2}, {7, 10}},
		},
		{
			"split interval",
			[]Interval[int]{{1, 10}},
			4, 6,
			[]Interval[int]{{1, 3}, {7, 10}},
		},
		{
			"multiple intervals mixed",
			[]Interval[int]{{1, 3}, {5, 8}, {10, 15}},
			6, 12,
			[]Interval[int]{{1, 3}, {5, 5}, {13, 15}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := s.Scalp(tt.xs, tt.lo, tt.hi)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestMergerInt(t *testing.T) {
	m := NewMergeImpl(idxInt, cmpInt)
	tests := []struct {
		name string
		xs0  []int
		xs1  []int
		want []int
	}{
		{
			"both empty",
			[]int{}, []int{}, []int{},
		},
		{
			"first empty",
			[]int{}, []int{1, 2, 3}, []int{1, 2, 3},
		},
		{
			"second empty",
			[]int{4, 5}, []int{}, []int{4, 5},
		},
		{
			"no overlap",
			[]int{1, 3, 5}, []int{2, 4, 6}, []int{1, 2, 3, 4, 5, 6},
		},
		{
			"all less",
			[]int{1, 2}, []int{3, 4}, []int{1, 2, 3, 4},
		},
		{
			"all greater",
			[]int{5, 6}, []int{1, 2}, []int{1, 2, 5, 6},
		},
		{
			"with duplicates",
			[]int{1, 2, 2, 3}, []int{2, 3, 4}, []int{1, 2, 2, 3, 4},
		},
		{
			"interleaved",
			[]int{1, 4, 7}, []int{2, 3, 5, 6}, []int{1, 2, 3, 4, 5, 6, 7},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := m.Merge(tt.xs0, tt.xs1)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestClamperInt(t *testing.T) {
	c := NewClampImpl(idxInt, cmpInt)

	tests := []struct {
		name string
		xs   []int
		lo   int
		hi   int
		want []int
	}{
		{
			"empty slice",
			[]int{}, 3, 7, []int{},
		},
		{
			"all below",
			[]int{1, 2}, 3, 5, []int{},
		},
		{
			"all above",
			[]int{10, 12}, 3, 8, []int{},
		},
		{
			"within range",
			[]int{3, 4, 5, 6, 7}, 3, 7, []int{3, 4, 5, 6, 7},
		},
		{
			"partial low",
			[]int{1, 3, 5}, 3, 6, []int{3, 5},
		},
		{
			"partial high",
			[]int{4, 6, 8}, 2, 6, []int{4, 6},
		},
		{
			"mixed",
			[]int{0, 2, 4, 6, 8}, 3, 7, []int{4, 6},
		},
		{
			"boundaries excluded below",
			[]int{2, 3, 4}, 3, 4, []int{3, 4},
		},
		{
			"boundaries excluded above",
			[]int{4, 5, 6}, 4, 5, []int{4, 5},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := c.Clamp(tt.xs, tt.lo, tt.hi)
			assert.Equal(t, tt.want, got)
		})
	}
}
