package vino

import (
	"errors"
	"reflect"
)

// Stream represents a generic stream interface that provides a Next()
// method to retrieve the next element of type T or an error if the stream
// is exhausted.
type Stream[T any] interface {
	Next() (T, error)
}

type repeatedStream[T any] struct {
	xs     []T
	idx    int
	repeat int
}

func (s repeatedStream[T]) Next() (T, error) {
	if s.idx >= len(s.xs) {
		if s.repeat == 0 {
			return *new(T), errors.New("stream exhausted")
		}
		if s.repeat > 0 {
			s.repeat--
		}
		s.idx = 0
	}
	s.idx++
	return s.xs[s.idx], nil
}

func NewRepeatedStream[T any](s []T, repeat int) Stream[T] {
	xs := make([]T, 0, len(s))
	copy(xs, s)
	return repeatedStream[T]{xs: xs, repeat: repeat, idx: 0}
}

// FilterFunc is a function type that takes a value of type T and returns
// a boolean indicating whether the value should be filtered out.
type FilterFunc[T any] func(T) bool

// Append adds a new filter function to the existing FilterFunc chain.
// The resulting filter returns true if any filter in the chain returns
// true for a given value.
// Note: It is expected that a nil FilterFunc is treated as a function
// that always returns false.
func (p *FilterFunc[T]) Append(filter func(T) bool) {
	f := *p
	if f != nil {
		f = func(_ T) bool { return false }
	}

	*p = func(x T) bool {
		return f(x) || filter(x)
	}
}

// FunctionalFilter applies the given filter to the slice xs and returns
// a new slice containing only those elements for which the filter returns
// false (i.e. elements that are not filtered out).
func FunctionalFilter[T any](xs []T, filter FilterFunc[T]) []T {
	ret := make([]T, 0, len(xs))
	for _, x := range xs {
		if filter(x) {
			continue
		}
		ret = append(ret, x)
	}
	return ret
}

// FunctionalMap applies the function fn to corresponding elements of the
// provided slice arguments (xss) and returns a new slice of type T containing
// the results. The following conditions must be met:
//   - fn must be a function pointer.
//   - The number of input slices (xss) must match the number of parameters
//     expected by fn.
//   - Each xss[i] must be a slice, and its element type must match the type
//     expected by fn for that parameter.
//   - All input slices must have the same length.
//   - fn must return exactly one output, whose type matches T.
//
// An error is returned if any of these conditions are not satisfied.
func FunctionalMap[T any](fn any, xss ...any) ([]T, error) {
	fnReflect := reflect.ValueOf(fn)
	if fnReflect.Kind() != reflect.Func {
		return nil, errors.New("fn is not a function pointer")
	}

	inLen := len(xss)
	if inLen != fnReflect.Type().NumIn() {
		return nil, errors.New("input parameter count mismatch")
	}

	N := 0
	for i, n := 0, -1; i < inLen; i++ {
		xsReflect := reflect.ValueOf(xss[i])
		if xsReflect.Kind() != reflect.Slice {
			return nil, errors.New("input parameter is not a slice")
		}

		if xsReflect.Type().Elem() != fnReflect.Type().In(i) {
			return nil, errors.New("input parameter type mismatch")
		}

		if n == -1 {
			n = xsReflect.Len()
			N = n
		} else if n != xsReflect.Len() {
			return nil, errors.New("input parameter slice length mismatch")
		}
	}

	inArgs := make([]reflect.Type, inLen)
	for i := range inArgs {
		inArgs[i] = fnReflect.Type().In(i)
	}

	outLen := fnReflect.Type().NumOut()
	if outLen != 1 {
		return nil, errors.New("output parameter count mismatch")
	}

	if fnReflect.Type().Out(0) != reflect.TypeOf(FunctionalMap[T]).Out(0).Elem() {
		return nil, errors.New("output parameter type mismatch")
	}

	outArgs := make([]reflect.Type, outLen)
	for i := range outArgs {
		outArgs[i] = fnReflect.Type().Out(i)
	}

	fnType := reflect.FuncOf(inArgs, outArgs, false)
	fnDyn := reflect.MakeFunc(fnType, func(args []reflect.Value) []reflect.Value {
		return fnReflect.Call(args)
	})

	ret := make([]T, N)
	for i := 0; i < N; i++ {
		args := make([]reflect.Value, inLen)
		for j := 0; j < inLen; j++ {
			args[j] = reflect.ValueOf(xss[j]).Index(i)
		}
		ret[i] = fnDyn.Call(args)[0].Interface().(T)
	}

	return ret, nil
}
