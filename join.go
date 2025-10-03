package rxgo

import "iter"

func StartWith[T any](values ...T) OperatorFunc[T, T] {
	return func(input Observable[T]) Observable[T] {
		return (ObservableFunc[T])(func(yield func(T, error) bool) {
			for len(values) > 0 {
				if !yield(values[0], nil) {
					return
				}
				values = values[1:]
			}

			next, stop := iter.Pull2((iter.Seq2[T, error])(input.Subscribe()))
			defer stop()

			for {
				v, err, ok := next()
				if err != nil {
					yield(v, err)
					return
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
