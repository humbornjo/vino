package vino

import "unsafe"

// ---------------------------------------------------------------------------
// Option Implementation
//
// This implementation provides a minimal "option" type using unsafe.Pointer.
// It encapsulates a pointer to a value of any type T and allows matching
// and extracting that value via a closure.
//
// The option type is defined as an alias for unsafe.Pointer. A nil option
// represents the absence of a value (i.e. None). When a non-nil pointer is
// provided, Option[T] creates a closure that, when invoked, assigns the
// stored value to a given target pointer. The Match function checks whether
// the option is None or contains a value and returns an extraction function.
//
// Note: This implementation uses unsafe operations. Caution is advised when
// integrating it into production code.
// ---------------------------------------------------------------------------

type some_f[T any] func(*T) option

var (
	// some_v is an auxiliary variable used to initialize the some value.
	some_v int = 1

	// None represents an option with no value.
	None option = option(nil)

	// `some` is a sentinel value representing the presence of a value. It is
	// computed using unsafe.Pointer conversion on some_v.
	some option = option(*(**unsafe.Pointer)(unsafe.Pointer(&some_v)))
)

// option is an alias for unsafe.Pointer and is used to represent
// an optional value.
type option unsafe.Pointer

// Option converts a pointer to a value of type T into an option like the one
// in Rust. If the provided pointer is nil, it returns None. Otherwise, it
// creates a closure that captures the value pointed to by x. When this
// closure is invoked (via Match), it copies the encapsulated value into the
// provided target pointer.
//
// Example:
//
//	x := 42
//	opt := Option(&x)
func Option[T any](x *T) option {
	if x == nil {
		return None
	}
	f := func(t *T) { *t = *x }
	return option(unsafe.Pointer(&f))
}

// Match inspects the provided option. If the option is None, it returns
// None and a nil extraction function. Otherwise, it returns a sentinel
// "some" value and a function that, when called with a pointer to T,
// invokes the stored closure to copy the encapsulated value into that
// pointer. This provides a mechanism to safely extract the stored value.
//
// Example usage:
//
//	val := new(T)
//	switch o, Some := Match\[T\](opt); o {
//	case None:
//	    // do things when None
//	case Some(val):
//	    // do things with retrieved val
//	}
func Match[T any](o option) (option, some_f[T]) {
	if o == None {
		return None, nil
	}
	return some, func(v *T) option {
		(*(*some_f[T])(o))(v)
		return some
	}
}
