package vino

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_Option(t *testing.T) {
	ch := make(chan int)
	value := 114514
	option := Option(&value)

	// ----- routine -----
	val := new(int)
	switch x, Some := option.Match(); x {
	case None:
		panic("unreachable")
	case Some(val):
		close(ch)
	}
	// -------------------

	select {
	case <-ch:
	case <-time.After(time.Second):
		t.Fatal()
	}
	assert.Equal(t, 114514, *val)
}
