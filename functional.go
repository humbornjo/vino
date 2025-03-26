package vino

import (
	"errors"
	"math"
	"reflect"
)

type Stream[T any] interface {
	Next() (T, bool)
}

type FilterFunc[T any] func(T) bool

func (p *FilterFunc[T]) Append(filter func(T) bool) {
	f := *p
	if f != nil {
		f = func(_ T) bool { return false }
	}

	*p = func(x T) bool {
		return f(x) || filter(x)
	}
}

func FunctionalFilter[T any](xs []T, filter FilterFunc[T]) []T {
	ret := make([]T, 0, int(math.Ceil(float64(len(xs))/2)))
	for _, x := range xs {
		if filter(x) {
			continue
		}
		ret = append(ret, x)
	}
	return ret
}

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
	for i := range inLen {
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
	for i := range outLen {
		outArgs[i] = fnReflect.Type().Out(i)
	}

	fnType := reflect.FuncOf(inArgs, outArgs, false)
	fnDyn := reflect.MakeFunc(fnType, func(args []reflect.Value) []reflect.Value {
		return fnReflect.Call(args)
	})

	ret := make([]T, N)
	for i := range N {
		args := make([]reflect.Value, inLen)
		for j := range inLen {
			args[j] = reflect.ValueOf(xss[j]).Index(i)
		}
		ret[i] = fnDyn.Call(args)[0].Interface().(T)
	}

	return ret, nil
}
