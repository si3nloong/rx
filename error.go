package rxgo

import (
	"iter"
)

func CatchError[I, O any](fn func(error) Observable[O]) OperatorFunc[I, Either[I, O]] {
	return func(input Observable[I]) Observable[Either[I, O]] {
		return func(yield func(Either[I, O], error) bool) {
			next, stop := iter.Pull2((iter.Seq2[I, error])(input))
			defer stop()

			for {
				v, err, ok := next()
				if err != nil {
					next2, stop2 := iter.Pull2((iter.Seq2[O, error])(fn(err)))
					defer stop2()

					for {
						v2, err2, ok2 := next2()
						if err2 != nil {
							if !yield(Either[I, O]{}, err2) {
								return
							}
						} else if !ok2 {
							return
						} else {
							if !yield(Either[I, O]{v: v2}, nil) {
								return
							}
						}
					}
				} else if !ok {
					return
				} else {
					if !yield(Either[I, O]{v: v}, nil) {
						return
					}
				}
			}
		}
	}
}

func Retry[T any](count uint) OperatorFunc[T, T] {
	return func(input Observable[T]) Observable[T] {
		return func(yield func(T, error) bool) {
			next, stop := iter.Pull2((iter.Seq2[T, error])(input))
			defer stop()

			var retryCount uint

			for {
				v, err, ok := next()
				if err != nil {
					// Retry if error
					if retryCount < count {
						stop()
						next, stop = iter.Pull2((iter.Seq2[T, error])(input))
						retryCount++
					} else {
						var zero T
						yield(zero, err)
						return
					}
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

func ThrowIfEmpty[T comparable](fn ...func() error) OperatorFunc[T, T] {
	return func(input Observable[T]) Observable[T] {
		return func(yield func(T, error) bool) {
			next, stop := iter.Pull2((iter.Seq2[T, error])(input))
			defer stop()

			var emptyValue T
			for {
				v, err, ok := next()
				if err != nil {
					yield(v, err)
					return
				} else if !ok {
					return
				} else {
					if v == emptyValue {
						if len(fn) > 0 {
							yield(emptyValue, fn[0]())
							return
						} else {
							yield(emptyValue, ErrEmpty)
							return
						}
					}
					if !yield(v, nil) {
						return
					}
				}
			}
		}
	}
}
