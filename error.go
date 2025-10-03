package rxgo

import (
	"iter"
)

func CatchError[T any](selector func(error) Observable[T]) OperatorFunc[T, T] {
	return func(input Observable[T]) Observable[T] {
		return (ObservableFunc[T])(func(yield func(T, error) bool) {
			next, stop := iter.Pull2(input.Subscribe())
			defer stop()

			for {
				v, err, ok := next()
				if err != nil {
					next2, stop2 := iter.Pull2(selector(err).Subscribe())
					defer stop2()

					for {
						v2, err2, ok2 := next2()
						if err2 != nil {
							var zero T
							if !yield(zero, err2) {
								return
							}
						} else if !ok2 {
							return
						} else {
							if !yield(v2, nil) {
								return
							}
						}
					}
				} else if !ok {
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

func CatchError2[I, O any](selector func(error) Observable[O]) OperatorFunc[I, Either[I, O]] {
	return func(input Observable[I]) Observable[Either[I, O]] {
		return (ObservableFunc[Either[I, O]])(func(yield func(Either[I, O], error) bool) {
			next, stop := iter.Pull2(input.Subscribe())
			defer stop()

			for {
				v, err, ok := next()
				if err != nil {
					next2, stop2 := iter.Pull2(selector(err).Subscribe())
					defer stop2()

					for {
						v2, err2, ok2 := next2()
						if err2 != nil {
							yield(Either[I, O]{}, err2)
							return
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
		})
	}
}

func Retry[T any](count int) OperatorFunc[T, T] {
	return func(input Observable[T]) Observable[T] {
		return (ObservableFunc[T])(func(yield func(T, error) bool) {
			next, stop := iter.Pull2(input.Subscribe())
			defer stop()

			var retryCount int
			for {
				v, err, ok := next()
				// Retry when it hit error
				if err != nil {
					if count < 0 /* Infinite retry */ {
						stop()
						next, stop = iter.Pull2(input.Subscribe())
					} else if retryCount < count {
						stop()
						next, stop = iter.Pull2(input.Subscribe())
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
		})
	}
}

func ThrowIfEmpty[T comparable](fn ...func() error) OperatorFunc[T, T] {
	return func(input Observable[T]) Observable[T] {
		return (ObservableFunc[T])(func(yield func(T, error) bool) {
			var emptyValue T
			for v, err := range input.Subscribe() {
				if err != nil {
					yield(v, err)
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
		})
	}
}
