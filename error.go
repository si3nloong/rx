package rxgo

import (
	"errors"
	"iter"
)

var ErrEmpty = errors.New(`rxgo: empty value`)
var ErrNotFound = errors.New(`rxgo: no values match`)

type Either[T1, T2 any] struct {
	v any
}

func (a Either[T1, T2]) Type1() (T1, bool) {
	if v, ok := a.v.(T1); ok {
		return v, true
	}
	var zero T1
	return zero, false
}

func (a Either[T1, T2]) Type2() (T2, bool) {
	if v, ok := a.v.(T2); ok {
		return v, true
	}
	var zero T2
	return zero, false
}

func (a Either[T1, T2]) MustType1() T1 {
	return a.v.(T1)
}
func (a Either[T1, T2]) MustType2() T2 {
	return a.v.(T2)
}

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

func ThrowError[T any](fn func() error) Observable[T] {
	return func(yield func(T, error) bool) {
		var v T
		yield(v, fn())
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
