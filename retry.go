package rxgo

import "iter"

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
