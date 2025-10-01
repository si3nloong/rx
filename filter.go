package rxgo

import "iter"

func Filter[T any](fn func(v T) bool) OperatorFunc[T, T] {
	return func(input Observable[T]) Observable[T] {
		return func(yield func(T, error) bool) {
			next, stop := iter.Pull2((iter.Seq2[T, error])(input))
			defer stop()

			for {
				v, err, ok := next()
				if err != nil {
					yield(v, err)
					return
				} else if !ok {
					return
				} else {
					if fn(v) {
						if !yield(v, nil) {
							return
						}
					}
				}
			}
		}
	}
}

func Filter2[T any](fn func(v T) (bool, error)) OperatorFunc[T, T] {
	return func(input Observable[T]) Observable[T] {
		return func(yield func(T, error) bool) {
			next, stop := iter.Pull2((iter.Seq2[T, error])(input))
			defer stop()

			for {
				v, err, ok := next()
				if err != nil {
					yield(v, err)
					return
				} else if !ok {
					return
				} else {
					ok, err := fn(v)
					if err != nil {
						yield(v, err)
						return
					}
					if ok {
						if !yield(v, nil) {
							return
						}
					}
				}
			}
		}
	}
}
