package vino

type option_t int
type some[T any] func(*T) option_t

var (
	None option_t = 0
	Some option_t = 1
)

type option[T any] struct {
	x     option_t
	somef func(*T)
}

func Option[T any](x *T) option[T] {
	if x == nil {
		return option[T]{None, func(t *T) {}}
	}

	return option[T]{Some, func(t *T) { *t = *x }}
}

func (o option[T]) Match() (option_t, some[T]) {
	return o.x, func(t *T) option_t {
		o.somef(t)
		return Some
	}
}
