package rx

import (
	"context"
	"iter"
	"sync"
	"time"
)

// Buffer buffers the source Observable values until closingNotifier emits.
func Buffer[T any, I any](closingNotifier Observable[I]) OperatorFunc[T, []T] {
	return func(input Observable[T]) Observable[[]T] {
		return (ObservableFunc[[]T])(func(yield func([]T, error) bool) {
			var buffer = make([]T, 0)
			var rw sync.RWMutex

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			go func() {
				next, stop := iter.Pull2(input.Subscribe())
				defer stop()

				for {
					select {
					case <-ctx.Done():
						return
					default:
						v, err, ok := next()
						if err != nil {
							yield(nil, err)
							return
						} else if !ok {
							return
						} else {
							rw.Lock()
							buffer = append(buffer, v)
							rw.Unlock()
						}
					}
				}
			}()

			for _, err := range closingNotifier.Subscribe() {
				if err != nil {
					cancel()
					yield(nil, err)
					return
				} else {
					rw.Lock()
					if !yield(buffer, nil) {
						rw.Unlock()
						return
					}
					buffer = make([]T, 0)
					rw.Unlock()
				}
			}
		})
	}
}

// BufferCount buffers the source Observable values into a slice of a specific size.
func BufferCount[T any](count uint) OperatorFunc[T, []T] {
	return func(input Observable[T]) Observable[[]T] {
		return (ObservableFunc[[]T])(func(yield func([]T, error) bool) {
			buffer := make([]T, 0, count)
			for v, err := range input.Subscribe() {
				if err != nil {
					yield(nil, err)
					return
				} else {
					buffer = append(buffer, v)
					if (uint)(len(buffer)) >= count {
						if !yield(buffer, nil) {
							return
						}
						buffer = make([]T, 0, count)
					}
				}
			}
			if len(buffer) > 0 {
				if !yield(buffer, nil) {
					return
				}
				clear(buffer)
			}
		})
	}
}

// Buffers the source Observable values for a specific time period.
func BufferTime[T any](duration time.Duration) OperatorFunc[T, []T] {
	return func(input Observable[T]) Observable[[]T] {
		return (ObservableFunc[[]T])(func(yield func([]T, error) bool) {
			next, stop := iter.Pull2(input.Subscribe())
			defer stop()

			timer := time.NewTicker(duration)
			defer timer.Stop()

			buffer := make([]T, 0)
			for {
				select {
				case <-timer.C:
					if !yield(buffer, nil) {
						return
					}
					buffer = make([]T, 0)
				default:
					v, err, ok := next()
					if err != nil {
						return
					} else if !ok {
						return
					} else {
						buffer = append(buffer, v)
					}
				}
			}
		})
	}
}

// Map applies a given project function to each value emitted by the source Observable, and emits the resulting values as an Observable.
func Map[I, O any](fn func(v I, index int) O) OperatorFunc[I, O] {
	return func(input Observable[I]) Observable[O] {
		return (ObservableFunc[O])(func(yield func(O, error) bool) {
			var i int
			for v, err := range input.Subscribe() {
				if err != nil {
					var o O
					yield(o, err)
					return
				} else {
					if !yield(fn(v, i), nil) {
						return
					}
				}
				i++
			}
		})
	}
}

// MapErr is similar to Map but deals with error.
func MapErr[I, O any](fn func(v I, index int) (O, error)) OperatorFunc[I, O] {
	return func(input Observable[I]) Observable[O] {
		return (ObservableFunc[O])(func(yield func(O, error) bool) {
			var i int
			for v, err := range input.Subscribe() {
				if err != nil {
					var zero O
					yield(zero, err)
					return
				} else {
					o, err := fn(v, i)
					if err != nil {
						yield(o, err)
						return
					}
					if !yield(o, nil) {
						return
					}
				}
				i++
			}
		})
	}
}

// ConcatMap projects each source value to an Observable which is merged in the output Observable, in a serialized fashion waiting for each one to complete before merging the next.
func ConcatMap[I, O any](project func(v I, index int) Observable[O]) OperatorFunc[I, O] {
	return func(input Observable[I]) Observable[O] {
		return (ObservableFunc[O])(func(yield func(O, error) bool) {
			var i int
			for v, err := range input.Subscribe() {
				if err != nil {
					var zero O
					yield(zero, err)
					return
				}
				for v2, err2 := range project(v, i).Subscribe() {
					if err2 != nil {
						var zero O
						yield(zero, err2)
						return
					} else {
						if !yield(v2, nil) {
							return
						}
					}
				}
				i++
			}
		})
	}
}

// SwitchMap projects each source value to an Observable which is merged in the output Observable, emitting values only from the most recently projected Observable.
func SwitchMap[I, O any](fn func(v I, index int) Observable[O]) OperatorFunc[I, O] {
	return func(input Observable[I]) Observable[O] {
		return (ObservableFunc[O])(func(yield func(O, error) bool) {
			next, stop := iter.Pull2(input.Subscribe())
			defer stop()

			var i int

			for {
				v, err, ok := next()
				if err != nil {
					var zero O
					yield(zero, err)
					return
				} else if !ok {
					return
				} else {
					next2, stop2 := iter.Pull2(fn(v, i).Subscribe())

				loop2:
					for {
						v2, err2, ok2 := next2()
						if err2 != nil {
							stop2()
							var zero O
							yield(zero, err2)
							return
						} else if !ok2 {
							break loop2
						} else {
							if !yield(v2, nil) {
								stop2()
								return
							}
						}
					}
					stop2()
				}
				i++
			}
		})
	}
}

// MergeMap projects each source value to an Observable which is merged in the output Observable.
func MergeMap[T any](fn func(v T, index int) Observable[T]) OperatorFunc[T, T] {
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
					next2, stop2 := iter.Pull2(fn(v, i).Subscribe())

				loop2:
					for {
						v2, err2, ok2 := next2()
						if err2 != nil {
							var zero T
							yield(zero, err2)
							return
						} else if !ok2 {
							break loop2
						} else {
							if !yield(v2, nil) {
								return
							}
						}
					}
					stop2()
				}
				i++
			}
		})
	}
}

// Pairwise groups pairs of consecutive emissions together and emits them as an array of two values.
func Pairwise[T any]() OperatorFunc[T, [2]T] {
	return func(input Observable[T]) Observable[[2]T] {
		return (ObservableFunc[[2]T])(func(yield func([2]T, error) bool) {
			pair := make([]T, 0, 2)
			for v, err := range input.Subscribe() {
				if err != nil {
					var zero [2]T
					yield(zero, err)
					return
				} else {
					pair = append(pair, v)
					if len(pair) > 1 {
						if !yield([2]T(pair), nil) {
							return
						}
						pair = pair[1:]
					}
				}
			}
		})
	}
}

// Scan applies an accumulator function over the source Observable, and returns each intermediate result, with the specified seed as the initial accumulator value.
func Scan[V, A any](accumulator func(acc A, value V, index int) A, seed A) OperatorFunc[V, A] {
	return func(input Observable[V]) Observable[A] {
		return (ObservableFunc[A])(func(yield func(A, error) bool) {
			var (
				acc = seed
				i   int
			)
			for v, err := range input.Subscribe() {
				if err != nil {
					var zero A
					yield(zero, err)
					return
				} else {
					acc = accumulator(acc, v, i)
					if !yield(acc, nil) {
						return
					}
					i++
				}
			}
		})
	}
}
