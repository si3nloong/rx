package rxgo

import (
	"context"
	"iter"
)

func Take[T any](count uint) OperatorFunc[T, T] {
	return func(input Observable[T]) Observable[T] {
		return (ObservableFunc[T])(func(yield func(T, error) bool) {
			next, stop := iter.Pull2((iter.Seq2[T, error])(input.Subscribe()))
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
		})
	}
}

func TakeLast[T any](count uint) OperatorFunc[T, T] {
	return func(input Observable[T]) Observable[T] {
		return (ObservableFunc[T])(func(yield func(T, error) bool) {
			next, stop := iter.Pull2((iter.Seq2[T, error])(input.Subscribe()))
			defer stop()

			result := make([]T, 0, count)
		loop:
			for {
				v, err, ok := next()
				if err != nil {
					var zero T
					yield(zero, err)
					return
				} else if !ok {
					break loop
				} else {
					result = append(result, v)
					if (uint)(len(result)) > count {
						result = result[1:]
					}
				}
			}

			for _, v := range result {
				if !yield(v, nil) {
					return
				}
			}
		})
	}
}

func TakeWhile[T any](fn func(T, int) bool) OperatorFunc[T, T] {
	return func(input Observable[T]) Observable[T] {
		return (ObservableFunc[T])(func(yield func(T, error) bool) {
			next, stop := iter.Pull2((iter.Seq2[T, error])(input.Subscribe()))
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
		})
	}
}

func TakeUntil[T, U any](notifier Observable[U]) OperatorFunc[T, T] {
	return func(input Observable[T]) Observable[T] {
		return (ObservableFunc[T])(func(yield func(T, error) bool) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			go func() {
				next, stop := iter.Pull2((iter.Seq2[U, error])(notifier.Subscribe()))
				defer stop()

				for {
					select {
					case <-ctx.Done():
						return
					default:
						if _, err, ok := next(); err != nil {
							return
						} else if !ok {
							return
						} else {
							cancel()
							return
						}
					}
				}
			}()

			next, stop := iter.Pull2((iter.Seq2[T, error])(input.Subscribe()))
			defer stop()

			for {
				select {
				case <-ctx.Done():
					return
				default:
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
		})
	}
}
