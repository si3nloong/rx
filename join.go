package rx

func StartWith[T any](values ...T) OperatorFunc[T, T] {
	return func(input Observable[T]) Observable[T] {
		return (ObservableFunc[T])(func(yield func(T, error) bool) {
			for len(values) > 0 {
				if !yield(values[0], nil) {
					return
				}
				values = values[1:]
			}

			for v, err := range input.Subscribe() {
				if err != nil {
					yield(v, err)
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
