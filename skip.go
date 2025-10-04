package rxgo

import (
	"context"
	"iter"
	"sync"
)

// Returns an Observable that skips the first count items emitted by the source Observable.
func Skip[T any](count uint) OperatorFunc[T, T] {
	return func(input Observable[T]) Observable[T] {
		return (ObservableFunc[T])(func(yield func(T, error) bool) {
			var skipCount uint
			for v, err := range input.Subscribe() {
				if err != nil {
					var zero T
					yield(zero, err)
					return
				} else {
					skipCount++
					if skipCount > count {
						if !yield(v, nil) {
							return
						}
					}
				}
			}
		})
	}
}

// Skip a specified number of values before the completion of an observable.
func SkipLast[T any](count uint) OperatorFunc[T, T] {
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
				}
			}

			if (uint)(len(result)) > count {
				result = result[:count]
				for _, v := range result {
					if !yield(v, nil) {
						return
					}
				}
			}
		})
	}
}

// Returns an Observable that skips all items emitted by the source Observable as long as a specified condition holds true, but emits all further source items as soon as the condition becomes false.
func SkipWhile[T any](fn func(T, int) bool) OperatorFunc[T, T] {
	return func(input Observable[T]) Observable[T] {
		return (ObservableFunc[T])(func(yield func(T, error) bool) {
			var i int
			for v, err := range input.Subscribe() {
				if err != nil {
					yield(v, err)
					return
				} else {
					if !fn(v, i) {
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

// Returns an Observable that skips items emitted by the source Observable until a second Observable emits an item.
func SkipUntil[T, U any](notifier Observable[U]) OperatorFunc[T, T] {
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
				// Internally, the skipUntil operator subscribes to the passed in notifier ObservableInput (which gets converted to an Observable)
				// in order to recognize the emission of its first value.
				v, err, ok := next()
				// When notifier emits next, the operator unsubscribes from it and starts emitting the values of the source observable until it completes or errors.
				select {
				case <-ctx.Done():
				case ch <- state[U]{v, err, ok}:
				}
			})

			next, stop := iter.Pull2(input.Subscribe())
			defer stop()

			var allowedEmit bool
		loop:
			for {
				select {
				case o := <-ch:
					// It will never let the source observable emit any values if the notifier completes or throws an error without emitting a value before.
					if o.err != nil {
						var zero T
						yield(zero, o.err)
						return
					} else if !o.ok {
						return
					} else {
						allowedEmit = true
					}
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
					} else if allowedEmit {
						if !yield(v, nil) {
							return
						}
					}
				}
			}

			wg.Wait()
		})
	}
}
