package vino_test

import (
	"sync"
	"testing"
	"time"

	. "github.com/humbornjo/vino"
)

func TestChanMut(t *testing.T) {
	size := uint(5)
	ch := NewChanMut[int](size)
	defer ch.Close()

	// Test initial size
	if ch.Cap() != int(size) {
		t.Errorf("expected initial size %d, got %d", size, ch.Cap())
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
	for ch.Cap() != int(newSize) {
		select {
		case <-due:
			t.Errorf("expected resized length %d, got %d", newSize, ch.Cap())
		default:
		}
	}
}

func TestChanBroadcast(t *testing.T) {
	const bufferSize = 3
	cb := NewChanBroadcast[int](bufferSize)

	// Register two receivers
	r1 := cb.Out()
	r2 := cb.Out()

	// Send data
	cb.In() <- 42
	cb.In() <- 100

	// Check if both receivers get the data
	select {
	case v := <-r1:
		if v != 42 {
			t.Errorf("Expected 42, got %d", v)
		}
	case <-time.After(time.Second):
		t.Errorf("Timeout waiting for message on r1")
	}

	select {
	case v := <-r2:
		if v != 42 {
			t.Errorf("Expected 42, got %d", v)
		}
	case <-time.After(time.Second):
		t.Errorf("Timeout waiting for message on r2")
	}

	// Close the broadcaster
	cb.Close()

	// Check if both receivers get the data
	select {
	case v := <-r1:
		if v != 100 {
			t.Errorf("Expected 100, got %d", v)
		}
	case <-time.After(time.Second):
		t.Errorf("Timeout waiting for message on r1")
	}

	select {
	case v := <-r2:
		if v != 100 {
			t.Errorf("Expected 100, got %d", v)
		}
	case <-time.After(time.Second):
		t.Errorf("Timeout waiting for message on r2")
	}

	// Ensure channels are closed
	_, ok := <-r1
	if ok {
		t.Errorf("Expected r1 to be closed")
	}
	_, ok = <-r2
	if ok {
		t.Errorf("Expected r2 to be closed")
	}
}

func TestChanBroadcast_MultiGuest(t *testing.T) {
	cb := NewChanBroadcast[int](5)

	r1, r2, r3 := cb.Out(), cb.Out(), cb.Out()

	wg := sync.WaitGroup{}
	N := 100

	// Send messages
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < N; i++ {
			cb.In() <- i
		}
	}()

	// Listen for messages
	checkReceived := func(r <-chan int, symbol string) {
		defer wg.Done()
		expected := 0
		for v := range r {
			if v != expected {
				t.Errorf("%s Expected %d, got %d", symbol, expected, v)
			}
			expected++
			if expected == N {
				break
			}
		}
	}

	wg.Add(3)
	go checkReceived(r1, "r1")
	go checkReceived(r2, "r2")
	go checkReceived(r3, "r3")

	wg.Wait()
	cb.Close()
}
