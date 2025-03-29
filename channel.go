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
	tunnel  chan T
	output  chan T
	resizer chan option
}

func NewChanMut[T any](size uint) *chanMut[T] {
	chanMut := &chanMut[T]{
		size:    size,
		input:   make(chan T),
		tunnel:  make(chan T, size),
		output:  make(chan T),
		resizer: make(chan option, 8),
	}
	go chanMut.prologue(size)
	go chanMut.epilogue()
	return chanMut
}

func (c *chanMut[T]) prologue(size uint) {
	defer func() {
		for range c.resizer {
		}
		close(c.resizer)
	}()

	resizing := false
	for {
		select {
		case x, ok := <-c.input:
			if !ok {
				close(c.tunnel)
				return
			}
			c.tunnel <- x
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
			newTunnel := make(chan T, size)
			c.tunnel, newTunnel = newTunnel, c.tunnel
			close(newTunnel)
			go func() {
				for x := range newTunnel {
					c.output <- x
				}
				resizing = false
				c.resizer <- None
			}()
		}
	}
}

func (c *chanMut[T]) epilogue() {
	defer close(c.output)
	update := false
	for {
		for x := range c.tunnel {
			update = false
			c.output <- x
		}
		if update {
			return
		} else {
			update = true
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
	return int(len(c.tunnel))
}

func (c *chanMut[T]) Cap() int {
	return int(cap(c.tunnel))
}

func (c *chanMut[T]) Resize(size uint) (err error) {
	defer func() {
		recover()
	}()
	err = errors.New("channel is closed")
	select {
	case _, ok := <-c.resizer:
		if !ok {
			return
		}
	default:
		select {
		case c.resizer <- Option(&size):
		default:
			return errors.New("resize failed")
		}
	}
	return nil
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
