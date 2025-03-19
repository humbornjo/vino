package vino

// TODO:
type MutChan[T any] struct {
	in   chan T
	out  chan T
	size int
}
