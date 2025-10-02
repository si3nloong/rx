package rxgo

import (
	"context"
	"iter"
	"time"
)

func Delay[T any](duration time.Duration) OperatorFunc[T, T] {
	return func(input Observable[T]) Observable[T] {
		return func(yield func(T, error) bool) {
			<-time.After(duration)

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

func Tap[T any](fn func(T)) OperatorFunc[T, T] {
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
					fn(v)
					if !yield(v, nil) {
						return
					}
				}
				i++
			}
		}
	}
}

func WithTimeInterval[T any]() OperatorFunc[T, TimeInterval[T]] {
	return func(input Observable[T]) Observable[TimeInterval[T]] {
		return func(yield func(TimeInterval[T], error) bool) {
			next, stop := iter.Pull2((iter.Seq2[T, error])(input))
			defer stop()

			startFrom := time.Now().UTC()

			for {
				v, err, ok := next()
				if err != nil {
					yield(TimeInterval[T]{}, err)
					return
				} else if !ok {
					break
				} else {
					if !yield(TimeInterval[T]{Interval: time.Since(startFrom), Value: v}, nil) {
						return
					}
				}
			}
		}
	}
}

func WithTimestamp[T any]() OperatorFunc[T, Timestamp[T]] {
	return func(input Observable[T]) Observable[Timestamp[T]] {
		return func(yield func(Timestamp[T], error) bool) {
			next, stop := iter.Pull2((iter.Seq2[T, error])(input))
			defer stop()

			for {
				v, err, ok := next()
				if err != nil {
					yield(Timestamp[T]{}, err)
					return
				} else if !ok {
					break
				} else {
					if !yield(Timestamp[T]{Time: time.Now().UTC(), Value: v}, nil) {
						return
					}
				}
			}
		}
	}
}

func Timeout[T any](duration time.Duration) OperatorFunc[T, T] {
	return func(input Observable[T]) Observable[T] {
		return func(yield func(T, error) bool) {
			next, stop := iter.Pull2((iter.Seq2[T, error])(input))
			defer stop()

			ctx, cancel := context.WithTimeout(context.Background(), duration)
			defer cancel()

			ch := make(chan state[T], 1)
			go func() {
				v, err, ok := next()
				select {
				case <-ctx.Done():
				case ch <- state[T]{v, err, ok}:
					cancel()
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
		}
	}
}

func ToSlice[T any]() OperatorFunc[T, []T] {
	return func(input Observable[T]) Observable[[]T] {
		return func(yield func([]T, error) bool) {
			next, stop := iter.Pull2((iter.Seq2[T, error])(input))
			defer stop()

			result := make([]T, 0)

			for {
				v, err, ok := next()
				if err != nil {
					yield(nil, err)
					return
				} else if !ok {
					break
				} else {
					result = append(result, v)
				}
			}

			yield(result, nil)
		}
	}
}
