package rx

import (
	"iter"
	"time"
)

// Defer creates an Observable that, on subscription, calls an Observable factory to make an Observable for each new Observer.
func Defer[T any](observableFactory func() Observable[T]) Observable[T] {
	return (ObservableFunc[T])(func(yield func(T, error) bool) {
		for v, err := range observableFactory().Subscribe() {
			if err != nil {
				var zero T
				yield(zero, err)
				return
			}
			if !yield(v, nil) {
				return
			}
		}
	})
}

// Empty creates an Observable that emits no items to the Observer and immediately completes.
func Empty[T any]() Observable[T] {
	return (ObservableFunc[T])(func(yield func(T, error) bool) {})
}

// Interval creates an Observable that emits a sequence of integers spaced by a given time interval.
func Interval(duration time.Duration) Observable[int] {
	return (ObservableFunc[int])(func(yield func(int, error) bool) {
		var i int

		for {
			<-time.After(duration)
			if !yield(i, nil) {
				return
			}
			i++
		}
	})
}

// From creates an Observable from an Array, an array-like object, an iterable object, or an Observable-like object.
//
// Example:
//
//	rx.From[int]([]int{1, 2, 3})
func From[T any, V Iterator[T]](items V) Observable[T] {
	return (ObservableFunc[T])(func(yield func(T, error) bool) {
		switch vi := any(items).(type) {
		case []T:
			for _, v := range vi {
				if !yield(v, nil) {
					return
				}
			}
		case chan T:
			for v := range vi {
				if !yield(v, nil) {
					return
				}
			}
		case <-chan T:
			for v := range vi {
				if !yield(v, nil) {
					return
				}
			}
		case iter.Seq[T]:
			for v := range vi {
				if !yield(v, nil) {
					return
				}
			}
		case iter.Seq2[T, error]:
			for v, err := range vi {
				if err != nil {
					var zero T
					yield(zero, err)
					return
				}
				if !yield(v, nil) {
					return
				}
			}
		default:
			panic("unreachable")
		}
	})
}

// FromChannel creates an Observable from a channel.
func FromChannel[T any, C interface {
	chan T | <-chan T
}](items C) Observable[T] {
	return (ObservableFunc[T])(func(yield func(T, error) bool) {
		for v := range items {
			if !yield(v, nil) {
				return
			}
		}
	})
}

// Of converts the arguments to an observable sequence.
//
// Example:
//
//	rx.Of(1, 2, 3)
func Of[T any](items ...T) Observable[T] {
	return (ObservableFunc[T])(func(yield func(T, error) bool) {
		for _, v := range items {
			if !yield(v, nil) {
				return
			}
		}
	})
}

// Range creates an Observable that emits a sequence of numbers within a specified range.
func Range[T Number](start, count T) Observable[T] {
	return (ObservableFunc[T])(func(yield func(T, error) bool) {
		for ; start <= count; start++ {
			if !yield(start, nil) {
				return
			}
		}
	})
}

// ThrowError creates an Observable that emits no items to the Observer and immediately emits an error notification.
func ThrowError[T any](errFactory func() error) Observable[T] {
	return (ObservableFunc[T])(func(yield func(T, error) bool) {
		var zero T
		yield(zero, errFactory())
	})
}

// Timer creates an Observable that starts emitting after an `initialDelay` and emits increasing numbers after each `period` of time thereafter.
func Timer[N Number](duration time.Duration) Observable[N] {
	return (ObservableFunc[N])(func(yield func(N, error) bool) {
		<-time.After(duration)
		var zero N
		yield(zero, nil)
	})
}

// Iif decies at subscription time which Observable will actually be subscribed.
func Iif[A, B any](condition func() bool, trueResult Observable[A], falseResult Observable[B]) Observable[Either[A, B]] {
	return (ObservableFunc[Either[A, B]])(func(yield func(Either[A, B], error) bool) {
		if condition() {
			for v, err := range trueResult.Subscribe() {
				if err != nil {
					yield(Either[A, B]{}, err)
					return
				}
				if !yield(Either[A, B]{v: v}, nil) {
					return
				}
			}
			return
		}

		for v, err := range falseResult.Subscribe() {
			if err != nil {
				yield(Either[A, B]{}, err)
				return
			}
			if !yield(Either[A, B]{v: v}, nil) {
				return
			}
		}
	})
}
