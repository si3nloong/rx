package rxgo

import (
	"iter"
)

func First[T any]() OperatorFunc[T, T] {
	return func(input Observable[T]) Observable[T] {
		return func(yield func(T, error) bool) {
			next, stop := iter.Pull2((iter.Seq2[T, error])(input))
			defer stop()

			v, err, ok := next()
			if err != nil {
				yield(v, err)
			} else if ok {
				yield(v, nil)
			}
		}
	}
}

func Last[T any]() OperatorFunc[T, T] {
	return func(input Observable[T]) Observable[T] {
		return func(yield func(T, error) bool) {
			next, stop := iter.Pull2((iter.Seq2[T, error])(input))
			defer stop()

			var latestValue T
			for {
				v, err, ok := next()
				if err != nil {
					yield(v, err)
					return
				} else if !ok {
					break
				} else {
					latestValue = v
				}
			}

			yield(latestValue, nil)
		}
	}
}
