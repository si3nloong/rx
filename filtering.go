package rxgo

import (
	"iter"
	"reflect"
	"time"
)

func AuditTime[T any](duration time.Duration) OperatorFunc[T, T] {
	return func(input Observable[T]) Observable[T] {
		return (ObservableFunc[T])(func(yield func(T, error) bool) {
			var latestValue T
			var stopped bool
			var timer *time.Timer
			defer func() {
				if timer != nil {
					timer.Stop()
				}
			}()
			for v, err := range input.Subscribe() {
				if err != nil {
					var zero T
					yield(zero, err)
					return
				} else if stopped {
					return
				} else {
					latestValue = v
					if timer == nil {
						timer = time.AfterFunc(duration, func() {
							timer = nil
							if !yield(latestValue, nil) {
								stopped = true
								return
							}
						})
					}
				}
			}
		})
	}
}

func DebounceTime[T any](duration time.Duration) OperatorFunc[T, T] {
	return func(input Observable[T]) Observable[T] {
		return (ObservableFunc[T])(func(yield func(T, error) bool) {
			var latestValue T
			var stopped bool
			var timer *time.Timer
			defer func() {
				if timer != nil {
					timer.Stop()
				}
			}()
			for v, err := range input.Subscribe() {
				if err != nil {
					var zero T
					yield(zero, err)
					return
				} else if stopped {
					return
				} else {
					if timer != nil {
						timer.Stop()
					}
					timer = time.AfterFunc(duration, func() {
						if !yield(latestValue, nil) {
							stopped = true
							return
						}
					})
					latestValue = v
				}
			}
		})
	}
}

func Distinct[T any, K comparable](keySelector func(value T) K) OperatorFunc[T, T] {
	return func(input Observable[T]) Observable[T] {
		return (ObservableFunc[T])(func(yield func(T, error) bool) {
			keyCache := make(map[K]struct{})
			defer clear(keyCache)
			results := make([]T, 0)
			for v, err := range input.Subscribe() {
				if err != nil {
					yield(v, err)
					return
				} else {
					key := keySelector(v)
					if _, ok := keyCache[key]; ok {
						continue
					}
					keyCache[key] = struct{}{}
					results = append(results, v)
				}
			}
			for len(results) > 0 {
				if !yield(results[0], nil) {
					return
				}
				results = results[1:]
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

func Last[T any]() OperatorFunc[T, T] {
	return func(input Observable[T]) Observable[T] {
		return (ObservableFunc[T])(func(yield func(T, error) bool) {
			var latestValue T
			for v, err := range input.Subscribe() {
				if err != nil {
					yield(v, err)
					return
				}
				latestValue = v
			}
			yield(latestValue, nil)
		})
	}
}

func SampleTime[T any](duration time.Duration) OperatorFunc[T, T] {
	return func(input Observable[T]) Observable[T] {
		return (ObservableFunc[T])(func(yield func(T, error) bool) {
		})
	}
}

func Single[T any](predicate func(T, int) bool) OperatorFunc[T, T] {
	return func(input Observable[T]) Observable[T] {
		return (ObservableFunc[T])(func(yield func(T, error) bool) {
			var i int
			var value *T
			for v, err := range input.Subscribe() {
				if err != nil {
					var zero T
					yield(zero, err)
					return
				}
				if predicate(v, i) {
					if value != nil {
						var zero T
						yield(zero, ErrSequence)
						return
					}
					value = &v
				}
				i++
			}
			var zero T
			if i > 0 {
				if value != nil {
					yield(*value, nil)
					return
				}
				yield(zero, ErrNotFound)
			} else {
				yield(zero, ErrEmpty)
			}
		})
	}
}

func ThrottleTime[T any](duration time.Duration) OperatorFunc[T, T] {
	return func(input Observable[T]) Observable[T] {
		return (ObservableFunc[T])(func(yield func(T, error) bool) {
			var latestValue T
			var emitted bool
			var timer *time.Timer
			defer func() {
				if timer != nil {
					timer.Stop()
				}
			}()
			for v, err := range input.Subscribe() {
				if err != nil {
					var zero T
					yield(zero, err)
					return
				} else {
					latestValue = v
					if !emitted {
						if !yield(latestValue, nil) {
							return
						}
						emitted = true
						timer = time.AfterFunc(duration, func() {
							emitted = false
						})
					}
				}
			}
		})
	}
}
