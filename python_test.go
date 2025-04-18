package vino_test

import (
	"testing"

	. "github.com/humbornjo/vino"
	"github.com/stretchr/testify/assert"
)

func TestBisectInt(t *testing.T) {
	b := NewBisector(func(a, b int) int {
		return a - b
	})

	assert.Equal(t, 0, b.BisectLeft([]int{1, 2, 3, 5}, 1))
	assert.Equal(t, 1, b.BisectLeft([]int{1, 2, 3, 5}, 2))
	assert.Equal(t, 2, b.BisectLeft([]int{1, 2, 3, 5}, 3))
	assert.Equal(t, 3, b.BisectLeft([]int{1, 2, 3, 5}, 4))
	assert.Equal(t, 3, b.BisectLeft([]int{1, 2, 3, 5}, 5))
	assert.Equal(t, 4, b.BisectLeft([]int{1, 2, 3, 5}, 6))

	assert.Equal(t, 1, b.BisectRight([]int{1, 2, 3, 5}, 1))
	assert.Equal(t, 2, b.BisectRight([]int{1, 2, 3, 5}, 2))
	assert.Equal(t, 3, b.BisectRight([]int{1, 2, 3, 5}, 3))
	assert.Equal(t, 3, b.BisectRight([]int{1, 2, 3, 5}, 4))
	assert.Equal(t, 4, b.BisectRight([]int{1, 2, 3, 5}, 5))
	assert.Equal(t, 4, b.BisectRight([]int{1, 2, 3, 5}, 6))
}
