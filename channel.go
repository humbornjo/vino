package vino

import (
	"errors"
	"sync"
	"weak"
)

// ------------------------------------------------------------------------
//  chanMut: Resizable Channel
// ------------------------------------------------------------------------

// chanMut is a resizable channel that supports dynamic adjustment of its
// buffer size. It maintains separate channels for input, internal
// processing (tunnel), and output, along with a resizer channel for
// handling resize requests.
type chanMut[T any] struct {
	size    uint
	input   chan T
	tunnel  chan T
	output  chan T
	resizer chan option
}

// NewChanMut creates and initializes a new resizable channel with the
// specified buffer size. The returned channel supports dynamic resizing
// via the Resize method.
func NewChanMut[T any](size uint) *chanMut[T] {
	chanMut := &chanMut[T]{
		size:    size,
		input:   make(chan T),
		tunnel:  make(chan T, size),
		output:  make(chan T),
		resizer: make(chan option, 8),
	}
	go chanMut.prologue()
	go chanMut.epilogue()
	return chanMut
}

func (c *chanMut[T]) prologue() {
	defer func() {
		close(c.resizer)
		for range c.resizer {
		}
	}()

	size := c.size
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

// Close gracefully shuts down the channel by closing the input channel.
// This will eventually close the underlying channels after processing
// pending messages.
func (c *chanMut[T]) Close() {
	close(c.input)
}

// In returns the send-only channel used to send messages into the
// resizable channel.
func (c *chanMut[T]) In() chan<- T {
	return c.input
}

// Out returns the receive-only channel used to receive messages from
// the resizable channel.
func (c *chanMut[T]) Out() <-chan T {
	return c.output
}

// Len returns the current number of messages buffered in channel.
func (c *chanMut[T]) Len() int {
	return int(len(c.tunnel))
}

// Cap returns the capacity of the channel.
func (c *chanMut[T]) Cap() int {
	return int(cap(c.tunnel))
}

// Resize attempts to change the buffer size of the underlying channel.
// It returns an error if the channel is closed or if resizing fails.
func (c *chanMut[T]) Resize(size uint) (err error) {
	defer func() { recover() }()
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

// ------------------------------------------------------------------------
//  chanBroadcast: Broadcast Channel
// ------------------------------------------------------------------------

// chanBroadcast is a channel that broadcasts messages to multiple
// registered receivers. It is recommended to use a for-range loop on
// the channel returned by Out() to ensure that no broadcast message
// is missed.
type chanBroadcast[T any] struct {
	closed    bool
	mu        sync.Mutex
	size      uint
	input     chan T
	tunnel    chan option
	registers map[weak.Pointer[chan T]][]T
}

// NewChanBroadcast creates and initializes a new broadcast channel with
// the specified buffer size. The broadcast channel distributes incoming
// messages to all registered receivers. User is responsible for setting
// the channel from Out() to nil when it is no longer needed so that GC
// can recycle the memory of channel and pending messages.
func NewChanBroadcast[T any](size uint) *chanBroadcast[T] {
	chanBroadcast := &chanBroadcast[T]{
		closed:    false,
		size:      size,
		input:     make(chan T),
		tunnel:    make(chan option, size),
		registers: make(map[weak.Pointer[chan T]][]T, 8),
	}
	go chanBroadcast.prologue()
	go chanBroadcast.epilogue()
	return chanBroadcast
}

func (c *chanBroadcast[T]) prologue() {
	for x := range c.input {
		c.tunnel <- Option(&x)
	}
	close(c.tunnel)
}

func (c *chanBroadcast[T]) epilogue() {
	defer func() {
		c.mu.Lock()
		defer c.mu.Unlock()
		c.closed = true

		for ch := range c.registers {
			if ptrCh := ch.Value(); ptrCh != nil {
				close(*ptrCh)
			}
			delete(c.registers, ch)
		}
	}()

	for x := range c.tunnel {
		v := new(T)
		dose := false
		encore := false

		switch o, Some := Match[T](x); o {
		case None:
		case Some(v):
			dose = true
		}

		c.mu.Lock()
		for ch, cargo := range c.registers {
			if ptrCh := ch.Value(); ptrCh != nil {
				if dose {
					cargo = append(cargo, *v)
				}
				if len(cargo) > 0 {
					select {
					case *ptrCh <- cargo[0]:
						cargo = cargo[1:]
					default:
						encore = true
					}
				}
				c.registers[ch] = cargo
			}
		}
		c.mu.Unlock()

		if encore {
			go func() {
				defer func() { recover() }()
				c.tunnel <- None
			}()
		}
	}
}

// Close shuts down the broadcast channel by closing its input channel.
// Once closed, all registered receiver channels will be closed. Pending
// messages will be dropped.
func (c *chanBroadcast[T]) Close() {
	close(c.input)
}

// In returns the send-only channel used for broadcasting messages.
func (c *chanBroadcast[T]) In() chan<- T {
	return c.input
}

// Out registers a new receiver for the broadcast channel and returns a
// receive-only channel. It is recommended to use a for-range loop to iterate
// over the returned channel if you don't want to miss any broadcast messages.
func (c *chanBroadcast[T]) Out() <-chan T {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.closed {
		return c.input
	}
	ch := make(chan T, c.size)
	c.registers[weak.Make(&ch)] = nil
	return ch
}

// Len returns the current number of broadcast messages in tunnel channel.
func (c *chanBroadcast[T]) Len() int {
	return len(c.tunnel)
}

// Cap returns the capacity of the broadcast channel's tunnel channel buffer.
func (c *chanBroadcast[T]) Cap() int {
	return cap(c.tunnel)
}
