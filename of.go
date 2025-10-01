package rxgo

import (
	"time"
)

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
