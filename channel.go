package vino

import (
	"sync"
	"weak"

	"github.com/google/uuid"
)

// TODO:
// A resizable channel
type ChanMut[T any] struct {
	in   chan T
	out  chan T
	size int
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
