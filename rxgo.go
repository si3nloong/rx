package rxgo

import "errors"

type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64
}

type OperatorFunc[I any, O any] func(Observable[I]) Observable[O]

var ErrEmpty = errors.New(`rxgo: empty value`)
var ErrNotFound = errors.New(`rxgo: no values match`)
var ErrTimeout = errors.New(`rxgo: timeout`)

type state[T any] struct {
	v   T
	err error
	ok  bool
}
