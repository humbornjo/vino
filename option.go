package vino

import "unsafe"

type option_t int
type some_f[T any] func(*T) option

var (
	some_v int    = 1
	None   option = option(nil)
	some   option = option(*(**unsafe.Pointer)(unsafe.Pointer(&some_v)))
)

type option unsafe.Pointer

func Option[T any](x *T) option {
	if x == nil {
		return option(nil)
	}
	f := func(t *T) { *t = *x }
	return option(unsafe.Pointer(&f))
}

func Match[T any](o option) (option, some_f[T]) {
	if o == None {
		return None, nil
	}
	return some, func(t *T) option {
		(*(*some_f[T])(o))(t)
		return some
	}
}
