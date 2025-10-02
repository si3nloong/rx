package rxgo

import (
	"context"
	"iter"
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

func TakeUntil[I, O any](notifier Observable[O]) OperatorFunc[I, I] {
	return func(input Observable[I]) Observable[I] {
		return func(yield func(I, error) bool) {
			next, stop := iter.Pull2((iter.Seq2[I, error])(input))
			defer stop()

			next2, stop2 := iter.Pull2((iter.Seq2[O, error])(notifier))
			defer stop2()

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			go func() {
				for {
					if _, err, ok := next2(); err != nil {
						return
					} else if !ok {
						return
					} else {
						cancel()
						break
					}
				}
			}()

			for {
				select {
				case <-ctx.Done():
					stop2()
					return
				default:
					v, err, ok := next()
					if err != nil {
						stop2()
						var zero I
						yield(zero, err)
						return
					} else if !ok {
						stop2()
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
}
