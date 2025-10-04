package rxgo

import (
	"errors"
	"iter"
)

type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64
}

type Observable[T any] interface {
	Subscribe() iter.Seq2[T, error]
	SubscribeOn(onNext func(T), onFailed func(error), onCompleted func())
}

type OperatorFunc[I any, O any] func(Observable[I]) Observable[O]

var ErrEmpty = errors.New(`rxgo: empty value`)
var ErrNotFound = errors.New(`rxgo: no values match`)
var ErrTimeout = errors.New(`rxgo: timeout`)
var ErrArgumentOutOfRange = errors.New(`rxgo: out of range`)
var ErrSequence = errors.New(`rxgo: too many values match`)

var errEmptyObservable = errors.New("rxgo: empty observable")

type state[T any] struct {
	v   T
	err error
	ok  bool
}

type goState[T any] struct {
	idx int
	v   T
	err error
	ok  bool
}
