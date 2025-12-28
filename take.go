package rx

import (
	"context"
	"iter"
	"sync"
)

// Emits only the first count values emitted by the source Observable.
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
			takeCount = 0
		})
	}
}

// Waits for the source to complete, then emits the last N values from the source, as specified by the count argument.
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
			result = nil
		})
	}
}

// Emits values emitted by the source Observable so long as each value satisfies the given predicate, and then completes as soon as this predicate is not satisfied.
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
			i = 0
		})
	}
}

// Emits the values emitted by the source Observable until a notifier Observable emits a value.
func TakeUntil[T, U any](notifier Observable[U]) OperatorFunc[T, T] {
	return func(input Observable[T]) Observable[T] {
		return (ObservableFunc[T])(func(yield func(T, error) bool) {
			var wg sync.WaitGroup
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			ch := make(chan state[U], 1)
			defer close(ch)

			wg.Go(func() {
				next, stop := iter.Pull2(notifier.Subscribe())
				defer stop()

				v, err, ok := next()
				select {
				case <-ctx.Done():
				case ch <- state[U]{v, err, ok}:
				}
			})

			next, stop := iter.Pull2(input.Subscribe())
			defer stop()
		loop:
			for {
				select {
				case o := <-ch:
					if o.err != nil {
						var zero T
						yield(zero, o.err)
						return
					} else if !o.ok {
						// If the notifier doesn't emit any value and completes then takeUntil will pass all values.
						continue
					}
					return
				default:
					v, err, ok := next()
					if err != nil {
						cancel()
						var zero T
						yield(zero, err)
						break loop
					} else if !ok {
						cancel()
						break loop
					} else {
						if !yield(v, nil) {
							cancel()
							break loop
						}
					}
				}
			}

			wg.Wait()
		})
	}
}
