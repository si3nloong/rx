package rxgo

import (
	"context"
	"iter"
	"sync"
)

func Take[T any](count uint) OperatorFunc[T, T] {
	return func(input Observable[T]) Observable[T] {
		return (ObservableFunc[T])(func(yield func(T, error) bool) {
			var takeCount uint
			for v, err := range input.Subscribe() {
				if err != nil {
					var zero T
					yield(zero, err)
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
			result := make([]T, 0, count)
			for v, err := range input.Subscribe() {
				if err != nil {
					var zero T
					yield(zero, err)
					return
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
			var i int
			for v, err := range input.Subscribe() {
				if err != nil {
					var zero T
					yield(zero, err)
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
			var wg sync.WaitGroup
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			wg.Go(func() {
				next, stop := iter.Pull2(notifier.Subscribe())
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
			})

			next, stop := iter.Pull2(input.Subscribe())
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
