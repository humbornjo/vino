package vino_test

import (
	"testing"
	"time"

	. "github.com/humbornjo/vino"
)

func TestChanMut(t *testing.T) {
	size := uint(5)
	ch := NewChanMut[int](size)
	defer ch.Close()

	// Test initial size
	if ch.Len() != int(size) {
		t.Errorf("expected initial size %d, got %d", size, ch.Len())
	}

	// Test sending and receiving
	go func() {
		for i := 0; i < int(size); i++ {
			ch.In() <- i
		}
	}()

	for i := 0; i < int(size); i++ {
		select {
		case val := <-ch.Out():
			if val != i {
				t.Errorf("expected %d, got %d", i, val)
			}
		case <-time.After(time.Second):
			t.Errorf("timed out waiting for value %d", i)
		}
	}

	// Test resizing
	newSize := uint(10)
	err := ch.Resize(newSize)
	if err != nil {
		t.Errorf("resize failed: %v", err)
	}

	due := time.After(time.Second)
	for ch.Len() != int(newSize) {
		select {
		case <-due:
			t.Errorf("expected resized length %d, got %d", newSize, ch.Len())
		default:
		}
	}
}
