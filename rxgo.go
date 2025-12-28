package rx

import (
	"errors"
	"iter"
)

// Number is a generic constraints that represents all numeric types in Go.
type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64
}

// Iterator represents a sequence of values that can be iterated over.
// It supports slices, channels, and iterators (iter.Seq and iter.Seq2).
type Iterator[T any] interface {
	[]T | chan T | <-chan T | iter.Seq[T] | iter.Seq2[T, error]
}

// Iterator2 represents a sequence of key-value pairs that can be iterated over.
// It supports slices, maps, channels, and iterators.
type Iterator2[T any, K comparable] interface {
	[]T | map[K]T | chan T | <-chan T | iter.Seq[T] | iter.Seq2[T, error]
}

// Observable is an interface that represents a stream of data.
// It provides methods to subscribe to the stream and receive values.
type Observable[T any] interface {
	// Subscribe returns an iterator that yields values and errors from the Observable.
	Subscribe() iter.Seq2[T, error]
	// SubscribeOn subscribes to the Observable and executes the provided callbacks for each event.
	SubscribeOn(onNext func(T), onFailed func(error), onCompleted func())
}

// OperatorFunc is a function that transforms an Observable of type I into an Observable of type O.
type OperatorFunc[I any, O any] func(Observable[I]) Observable[O]

// ErrEmpty is returned when an Observable is empty but a value was expected.
var ErrEmpty = errors.New(`rxgo: empty value`)

// ErrNotFound is returned when a value matches a criteria is not found.
var ErrNotFound = errors.New(`rxgo: no values match`)

// ErrTimeout is returned when an operation times out.
var ErrTimeout = errors.New(`rxgo: timeout`)

// ErrArgumentOutOfRange is returned when an argument is out of the valid range.
var ErrArgumentOutOfRange = errors.New(`rxgo: out of range`)

// ErrSequence is returned when a sequence contains too many values (e.g. Single operator).
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
