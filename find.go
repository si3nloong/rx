package rxgo

import "iter"

func Find[T any](predicate func(T, int) bool) OperatorFunc[T, T] {
	return func(input Observable[T]) Observable[T] {
		return func(yield func(T, error) bool) {
			next, stop := iter.Pull2((iter.Seq2[T, error])(input))
			defer stop()

			var i int
			for {
				v, err, ok := next()
				if err != nil {
					yield(v, err)
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
