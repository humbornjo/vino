package vino

type Stream[T any] interface {
	Next() (T, bool)
}

// TODO:
func Map(fn any, xss ...any) {

}
