package rxgo

import "iter"

func ToSlice[T any]() OperatorFunc[T, []T] {
	return func(input Observable[T]) Observable[[]T] {
		return func(yield func([]T, error) bool) {
			next, stop := iter.Pull2((iter.Seq2[T, error])(input))
			defer stop()

			result := make([]T, 0)

			for {
				v, err, ok := next()
				if err != nil {
					yield(nil, err)
					return
				} else if !ok {
					break
				} else {
					result = append(result, v)
				}
			}

			yield(result, nil)
		}
	}
}
