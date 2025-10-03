package rxgo

import (
	"iter"
	"time"
)

func Of[T any](items []T) Observable[T] {
	return (ObservableFunc[T])(func(yield func(T, error) bool) {
		for _, v := range items {
			if !yield(v, nil) {
				return
			}
		}
	})
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

func Empty[T any]() Observable[T] {
	return (ObservableFunc[T])(func(yield func(T, error) bool) {})
}

func ThrowError[T any](fn func() error) Observable[T] {
	return (ObservableFunc[T])(func(yield func(T, error) bool) {
		var v T
		yield(v, fn())
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
			next, stop := iter.Pull2(trueResult.Subscribe())
			defer stop()

			for {
				v, err, ok := next()
				if err != nil {
					yield(Either[A, B]{}, err)
					return
				} else if !ok {
					return
				} else {
					if !yield(Either[A, B]{v}, nil) {
						return
					}
				}
			}
		}

		next, stop := iter.Pull2(falseResult.Subscribe())
		defer stop()

		for {
			v, err, ok := next()
			if err != nil {
				yield(Either[A, B]{}, err)
				return
			} else if !ok {
				return
			} else {
				if !yield(Either[A, B]{v}, nil) {
					return
				}
			}
		}
	})
}
