package rxgo

import "iter"

func DefaultIfEmpty[T any](defaultValue T) OperatorFunc[T, T] {
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
					if i > 0 {
						return
					}
					yield(defaultValue, nil)
					return
				} else {
					if !yield(v, nil) {
						return
					}
					i++
				}
			}
		}
	}
}

func Every[T any](predicate func(T, int) bool) OperatorFunc[T, bool] {
	return func(input Observable[T]) Observable[bool] {
		return func(yield func(bool, error) bool) {
			next, stop := iter.Pull2((iter.Seq2[T, error])(input))
			defer stop()

			var i int
			var passed = true
			for {
				v, err, ok := next()
				if err != nil {
					yield(false, err)
					return
				} else if !ok {
					yield(passed, nil)
					return
				} else {
					passed = passed && predicate(v, i)
				}
				i++
			}
		}
	}
}

func Find[T any](predicate func(T, int) bool) OperatorFunc[T, T] {
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
					var zero T
					yield(zero, ErrNotFound)
					return
				} else {
					if predicate(v, i) {
						yield(v, nil)
						return
					}
				}
				i++
			}
		}
	}
}

func FindIndex[T any](predicate func(T, int) bool) OperatorFunc[T, int] {
	return func(input Observable[T]) Observable[int] {
		return func(yield func(int, error) bool) {
			next, stop := iter.Pull2((iter.Seq2[T, error])(input))
			defer stop()

			var i int
			for {
				v, err, ok := next()
				if err != nil {
					yield(-1, err)
					return
				} else if !ok {
					yield(-1, nil)
					return
				} else {
					if predicate(v, i) {
						yield(i, nil)
						return
					}
				}
				i++
			}
		}
	}
}

func IsEmpty[T any]() OperatorFunc[T, bool] {
	return func(input Observable[T]) Observable[bool] {
		return func(yield func(bool, error) bool) {
			next, stop := iter.Pull2((iter.Seq2[T, error])(input))
			defer stop()

			var i int
			for {
				if _, err, ok := next(); err != nil {
					yield(false, err)
					return
				} else if !ok {
					println("HERE OK", i)
					if i > 0 {
						yield(false, nil)
					} else {
						yield(true, nil)
					}
					return
				} else {
					i++
				}
			}
		}
	}
}
