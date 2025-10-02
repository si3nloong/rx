package rxgo

import "time"

func Of[T any](items []T) Observable[T] {
	return func(yield func(T, error) bool) {
		for _, v := range items {
			if !yield(v, nil) {
				return
			}
		}
	}
}

func Interval(duration time.Duration) Observable[int] {
	return func(yield func(int, error) bool) {
		var i int

		for {
			<-time.After(duration)
			if !yield(i, nil) {
				return
			}
			i++
		}
	}
}

func Empty[T any]() Observable[T] {
	return func(yield func(T, error) bool) {}
}

func ThrowError[T any](fn func() error) Observable[T] {
	return func(yield func(T, error) bool) {
		var v T
		yield(v, fn())
	}
}

func Timer(duration time.Duration) Observable[int] {
	return func(yield func(int, error) bool) {
		<-time.After(duration)
		yield(0, nil)
	}
}
