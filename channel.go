package vino

import (
	"errors"
	"sync"
	"weak"

	"github.com/google/uuid"
)

// A resizable channel
type chanMut[T any] struct {
	size    uint
	input   chan T
	output  chan T
	resizer chan option
}

func NewChanMut[T any](size uint) *chanMut[T] {
	ch := &chanMut[T]{
		size:    size,
		input:   make(chan T, size),
		output:  make(chan T),
		resizer: make(chan option, 8),
	}
	go ch.start(size)
	return ch
}

func (c *chanMut[T]) start(size uint) {
	resizing := false
	for {
		select {
		case x, ok := <-c.input:
			if !ok {
				close(c.output)
				close(c.resizer)
				return
			}
			c.output <- x
		case osize := <-c.resizer:
			newSize := new(uint)
			switch o, Some := Match[uint](osize); o {
			case None:
			case Some(newSize):
				size = *newSize
			}
			if resizing || c.size == size {
				continue
			}
			resizing = true
			c.size = size
			newInput := make(chan T, size)
			c.input, newInput = newInput, c.input
			go func() {
				for x := range newInput {
					c.input <- x
				}
				resizing = false
				c.resizer <- None
			}()
		}
	}
}

func (c *chanMut[T]) Close() {
	close(c.input)
}

func (c *chanMut[T]) In() chan<- T {
	return c.input
}

func (c *chanMut[T]) Out() <-chan T {
	return c.output
}

func (c *chanMut[T]) Len() int {
	return int(c.size)
}

func (c *chanMut[T]) Resize(size uint) error {
	select {
	case c.resizer <- Option(&size):
		return nil
	default:
		return errors.New("resize failed")
	}
}

// TODO:
// A channel that drop all messages that stay in the buffer for too long
type ChanLost[T any] struct {
}

// TODO:
// A channel that broadcast messages
type ChanBroadcast[T any] struct {
	mu   sync.Mutex
	wmap map[uuid.UUID]weak.Pointer[T]
}

// func (c *ChanBroadcast[T]) Close() {
// }
//
// func (c *ChanBroadcast[T]) In(x T) {
// }
//
// func (c *ChanBroadcast[T]) Out() T {
// }
