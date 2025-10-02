package rxgo

import (
	"iter"
	"time"
)

func Delay[T any](duration time.Duration) OperatorFunc[T, T] {
	return func(input Observable[T]) Observable[T] {
		return func(yield func(T, error) bool) {
			t := time.NewTimer(duration)
			defer t.Stop()

			<-t.C
			t.Stop()

			next, stop := iter.Pull2((iter.Seq2[T, error])(input))
			defer stop()

			for {
				v, err, ok := next()
				if err != nil {
					var zero T
					yield(zero, err)
					return
				} else if !ok {
					return
				} else {
					if !yield(v, nil) {
						return
					}
				}
			}
		}
	}
}

func Tap[T any](fn func(T)) OperatorFunc[T, T] {
	return func(input Observable[T]) Observable[T] {
		return func(yield func(T, error) bool) {
			next, stop := iter.Pull2((iter.Seq2[T, error])(input))
			defer stop()

			var i int
			for {
				v, err, ok := next()
				if err != nil {
					var zero T
					yield(zero, err)
					return
				} else if !ok {
					return
				} else {
					fn(v)
					if !yield(v, nil) {
						return
					}
				}
				i++
			}
		}
	}
}

func ToSlice[T any]() OperatorFunc[T, []T] {
	return func(input Observable[T]) Observable[[]T] {
		return func(yield func([]T, error) bool) {
			next, stop := iter.Pull2((iter.Seq2[T, error])(input))
			defer stop()

			result := make([]T, 0)

			for {
				v, err, ok := next()
				if err != nil {
					yield(nil, err)
					return
				} else if !ok {
					break
				} else {
					result = append(result, v)
				}
			}

			yield(result, nil)
		}
	}
}
