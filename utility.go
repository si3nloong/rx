package rx

import (
	"context"
	"iter"
	"time"
)

// Tap performs a side effect for every emission on the source Observable, but returns an Observable that is identical to the source.
func Tap[T any](fn func(T)) OperatorFunc[T, T] {
	return func(input Observable[T]) Observable[T] {
		return (ObservableFunc[T])(func(yield func(T, error) bool) {
			for v, err := range input.Subscribe() {
				if err != nil {
					yield(v, err)
					return
				} else {
					fn(v)
					if !yield(v, nil) {
						return
					}
				}
			}
		})
	}
}

// Delay delays the emissions of items from the source Observable by a given timeout or until a given Date.
func Delay[T any](duration time.Duration) OperatorFunc[T, T] {
	return func(input Observable[T]) Observable[T] {
		return (ObservableFunc[T])(func(yield func(T, error) bool) {
			<-time.After(duration)

			for v, err := range input.Subscribe() {
				if err != nil {
					yield(v, err)
					return
				} else {
					if !yield(v, nil) {
						return
					}
				}
			}
		})
	}
}

// DelayWhen delays the emission of items from the source Observable by a given time span determined by the emissions of another Observable.
func DelayWhen[T, R any](delayDurationSelector func(value T, index int) Observable[R]) OperatorFunc[T, T] {
	return func(input Observable[T]) Observable[T] {
		return (ObservableFunc[T])(func(yield func(T, error) bool) {
			next, stop := iter.Pull2(input.Subscribe())
			defer stop()

			var i int
			for {
				v, err, ok := next()
				if err != nil {
					yield(v, err)
					return
				} else if !ok {
					return
				} else {
					delayDurationSelector(v, i).Subscribe()
					if !yield(v, nil) {
						return
					}
					i++
				}
			}
		})
	}
}

// WithTimeInterval adds the time interval since the last emission to the emitted value.
func WithTimeInterval[T any]() OperatorFunc[T, TimeInterval[T]] {
	return func(input Observable[T]) Observable[TimeInterval[T]] {
		return (ObservableFunc[TimeInterval[T]])(func(yield func(TimeInterval[T], error) bool) {
			startFrom := time.Now().UTC()
			for v, err := range input.Subscribe() {
				if err != nil {
					yield(TimeInterval[T]{}, err)
					return
				} else {
					if !yield(TimeInterval[T]{Interval: time.Since(startFrom), Value: v}, nil) {
						return
					}
				}
			}
		})
	}
}

// WithTimestamp attaches a timestamp to each item emitted by an observable indicating when it was emitted.
func WithTimestamp[T any]() OperatorFunc[T, Timestamp[T]] {
	return func(input Observable[T]) Observable[Timestamp[T]] {
		return (ObservableFunc[Timestamp[T]])(func(yield func(Timestamp[T], error) bool) {
			for v, err := range input.Subscribe() {
				if err != nil {
					yield(Timestamp[T]{}, err)
					return
				} else {
					if !yield(Timestamp[T]{Time: time.Now().UTC(), Value: v}, nil) {
						return
					}
				}
			}
		})
	}
}

// Timeout errors if the Observable does not emit a value within a specified time.
func Timeout[T any](duration time.Duration) OperatorFunc[T, T] {
	return func(input Observable[T]) Observable[T] {
		return (ObservableFunc[T])(func(yield func(T, error) bool) {
			ctx, cancel := context.WithTimeout(context.Background(), duration)
			defer cancel()

			next, stop := iter.Pull2(input.Subscribe())
			defer stop()

			ch := make(chan state[T], 1)
			defer close(ch)

			go func() {
				v, err, ok := next()
				select {
				case <-ctx.Done():
					return
				case ch <- state[T]{v, err, ok}:
					cancel()
					if err != nil || !ok {
						return
					}
				}
			}()

			select {
			case <-ctx.Done():
				var zero T
				yield(zero, ErrTimeout)
				return

			case r := <-ch:
				if r.err != nil {
					var zero T
					yield(zero, r.err)
					return
				} else if !r.ok {
					return
				} else {
					if !yield(r.v, nil) {
						return
					}
				}

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
		})
	}
}

// ToSlice collects all values from the source Observable into a slice.
func ToSlice[T any]() OperatorFunc[T, []T] {
	return func(input Observable[T]) Observable[[]T] {
		return (ObservableFunc[[]T])(func(yield func([]T, error) bool) {
			result := make([]T, 0)
			for v, err := range input.Subscribe() {
				if err != nil {
					yield(nil, err)
					return
				} else {
					result = append(result, v)
				}
			}
			yield(result, nil)
		})
	}
}
