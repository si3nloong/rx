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

func Interval(d time.Duration) Observable[int] {
	return func(yield func(int, error) bool) {
		var i int
		for {
			<-time.After(d)
			if !yield(i, nil) {
				return
			}
			i++
		}
	}
}

func Empty() Observable[any] {
	return func(yield func(any, error) bool) {}
}

func ThrowError[T any](fn func() error) Observable[T] {
	return func(yield func(T, error) bool) {
		var v T
		yield(v, fn())
	}
}

func Timer(duration time.Duration) Observable[int] {
	return func(yield func(int, error) bool) {
		t := time.NewTimer(duration)
		defer t.Stop()

		var i int
		for {
			<-t.C
			if !yield(i, nil) {
				return
			}
			i++
		}
	}
}
