package rxgo

import (
	"time"
)

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

func Empty[T any]() Observable[T] {
	return (ObservableFunc[T])(func(yield func(T, error) bool) {})
}

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

func Of[T any](items []T) Observable[T] {
	return (ObservableFunc[T])(func(yield func(T, error) bool) {
		for _, v := range items {
			if !yield(v, nil) {
				return
			}
		}
	})
}

func Range[T Number](start, count T) Observable[T] {
	return (ObservableFunc[T])(func(yield func(T, error) bool) {
		for ; start <= count; start++ {
			if !yield(start, nil) {
				return
			}
		}
	})
}

func ThrowError[T any](errFactory func() error) Observable[T] {
	return (ObservableFunc[T])(func(yield func(T, error) bool) {
		var zero T
		yield(zero, errFactory())
	})
}

func Timer[N Number](duration time.Duration) Observable[N] {
	return (ObservableFunc[N])(func(yield func(N, error) bool) {
		<-time.After(duration)
		var zero N
		yield(zero, nil)
	})
}

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
