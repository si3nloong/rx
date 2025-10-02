package rxgo

import (
	"iter"
)

func Concat[T any](inputs ...Observable[T]) Observable[T] {
	if len(inputs) < 2 {
		panic(`Concat required at least 2 observable`)
	}
	return func(yield func(T, error) bool) {
		for len(inputs) > 0 {
			next, stop := iter.Pull2((iter.Seq2[T, error])(inputs[0]))

			for {
				v, err, ok := next()
				if err != nil {
					stop()
					var zero T
					yield(zero, err)
					return
				} else if !ok {
					stop()
					break
				} else {
					if !yield(v, nil) {
						stop()
						return
					}
				}
			}
			stop()

			inputs = inputs[1:]
		}
	}
}
