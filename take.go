package rxgo

import (
	"iter"
	"time"
)

func Take[T any](count uint) OperatorFunc[T, T] {
	return func(input Observable[T]) Observable[T] {
		return func(yield func(T, error) bool) {
			next, stop := iter.Pull2((iter.Seq2[T, error])(input))
			defer stop()

			var takeCount uint
			for {
				v, err, ok := next()
				if err != nil {
					var zero T
					yield(zero, err)
					return
				} else if !ok {
					return
				} else {
					takeCount++
					if takeCount > count {
						return
					}
					if !yield(v, nil) {
						return
					}
				}
			}
		}
	}
}

func TakeLast[T any](count uint) OperatorFunc[T, T] {
	return func(input Observable[T]) Observable[T] {
		return func(yield func(T, error) bool) {
			next, stop := iter.Pull2((iter.Seq2[T, error])(input))
			defer stop()

			next()
		}
	}
}

func TakeWhile[T any](fn func(T, int) bool) OperatorFunc[T, T] {
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
					if fn(v, i) {
						if !yield(v, nil) {
							return
						}
					}
				}
				i++
			}
		}
	}
}

func TakeUntil[T any](d time.Duration) OperatorFunc[T, T] {
	return func(input Observable[T]) Observable[T] {
		return func(yield func(T, error) bool) {
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
