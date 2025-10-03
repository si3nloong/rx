package rxgo

import (
	"context"
	"iter"
	"reflect"
	"time"
)

func DebounceTime[T any](duration time.Duration) OperatorFunc[T, T] {
	return func(input Observable[T]) Observable[T] {
		return (ObservableFunc[T])(func(yield func(T, error) bool) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			ch := make(chan state[T], 1)
			defer close(ch)

			go func() {
				next, stop := iter.Pull2(input.Subscribe())
				defer stop()

				for {
					v, err, ok := next()
					select {
					case <-ctx.Done():
						return
					case ch <- state[T]{0, v, err, ok}:
						if err != nil || !ok {
							return
						}
					}
				}
			}()

			var latestValue T
			var emitted bool
			for {
				select {
				case <-time.After(duration):
					if emitted {
						if !yield(latestValue, nil) {
							return
						}
						emitted = false
					}
				case r := <-ch:
					if r.err != nil {
						cancel()
						var zero T
						yield(zero, r.err)
						return
					} else if !r.ok {
						cancel()
						return
					} else {
						latestValue = r.v
						emitted = true
					}
				}
			}
		})
	}
}

func DistinctUntilChanged[T any](comparator ...func(prev, curr T) bool) OperatorFunc[T, T] {
	return func(input Observable[T]) Observable[T] {
		return (ObservableFunc[T])(func(yield func(T, error) bool) {
			next, stop := iter.Pull2(input.Subscribe())
			defer stop()

			latestValue, err, ok := next()
			if err != nil {
				var zero T
				yield(zero, err)
				return
			} else if !ok {
				yield(latestValue, nil)
				return
			} else {
				if !yield(latestValue, nil) {
					return
				}
			}

			fn := func(prev, curr T) bool {
				return reflect.DeepEqual(prev, curr)
			}
			if len(comparator) > 0 {
				fn = comparator[0]
			}

			for {
				v, err, ok := next()
				if err != nil {
					yield(v, err)
					return
				} else if !ok {
					return
				} else {
					if fn(latestValue, v) {
						continue
					}
					if !yield(v, nil) {
						return
					}
					latestValue = v
				}
			}
		})
	}
}

func ElementAt[T any](index int, defaultValue ...T) OperatorFunc[T, T] {
	return func(input Observable[T]) Observable[T] {
		return (ObservableFunc[T])(func(yield func(T, error) bool) {
			var i int
			for v, err := range input.Subscribe() {
				if err != nil {
					var zero T
					yield(zero, err)
					return
				} else {
					if i == index {
						yield(v, nil)
						return
					}
				}
			}
			if len(defaultValue) > 0 {
				yield(defaultValue[0], nil)
			} else {
				var zero T
				yield(zero, ErrArgumentOutOfRange)
			}
		})
	}
}

func Filter[T any](fn func(v T) bool) OperatorFunc[T, T] {
	return func(input Observable[T]) Observable[T] {
		return (ObservableFunc[T])(func(yield func(T, error) bool) {
			for v, err := range input.Subscribe() {
				if err != nil {
					yield(v, err)
					return
				} else {
					if fn(v) {
						if !yield(v, nil) {
							return
						}
					}
				}
			}
		})
	}
}

func Filter2[T any](fn func(v T) (bool, error)) OperatorFunc[T, T] {
	return func(input Observable[T]) Observable[T] {
		return (ObservableFunc[T])(func(yield func(T, error) bool) {
			for v, err := range input.Subscribe() {
				if err != nil {
					yield(v, err)
					return
				} else {
					if ok, err := fn(v); err != nil {
						var zero T
						yield(zero, err)
						return
					} else if ok {
						if !yield(v, nil) {
							return
						}
					}
				}
			}
		})
	}
}

func First[T any]() OperatorFunc[T, T] {
	return func(input Observable[T]) Observable[T] {
		return (ObservableFunc[T])(func(yield func(T, error) bool) {
			for v, err := range input.Subscribe() {
				if err != nil {
					var zero T
					yield(zero, err)
					return
				} else {
					yield(v, nil)
					break
				}
			}
		})
	}
}

func Last[T any]() OperatorFunc[T, T] {
	return func(input Observable[T]) Observable[T] {
		return (ObservableFunc[T])(func(yield func(T, error) bool) {
			next, stop := iter.Pull2(input.Subscribe())
			defer stop()

			var latestValue T
		loop:
			for {
				v, err, ok := next()
				if err != nil {
					yield(v, err)
					return
				} else if !ok {
					break loop
				} else {
					latestValue = v
				}
			}

			yield(latestValue, nil)
		})
	}
}

func IgnoreElements[T any]() OperatorFunc[T, T] {
	return func(input Observable[T]) Observable[T] {
		return (ObservableFunc[T])(func(yield func(T, error) bool) {
			for _, err := range input.Subscribe() {
				if err != nil {
					var zero T
					yield(zero, err)
					return
				}
			}
		})
	}
}

func Single[T any](fn func(T, int) bool) OperatorFunc[T, T] {
	return func(input Observable[T]) Observable[T] {
		return (ObservableFunc[T])(func(yield func(T, error) bool) {
			next, stop := iter.Pull2(input.Subscribe())
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
