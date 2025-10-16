package rx

func DefaultIfEmpty[T any](defaultValue T) OperatorFunc[T, T] {
	return func(input Observable[T]) Observable[T] {
		return (ObservableFunc[T])(func(yield func(T, error) bool) {
			var emitted bool
			for v, err := range input.Subscribe() {
				if err != nil {
					var zero T
					yield(zero, err)
					return
				} else {
					if !yield(v, nil) {
						return
					}
					emitted = true
				}
			}
			if emitted {
				return
			}
			yield(defaultValue, nil)
		})
	}
}

func Every[T any](predicate func(T, int) bool) OperatorFunc[T, bool] {
	return func(input Observable[T]) Observable[bool] {
		return (ObservableFunc[bool])(func(yield func(bool, error) bool) {
			var i int
			var passed = true
			for v, err := range input.Subscribe() {
				if err != nil {
					yield(false, err)
					return
				}
				passed = passed && predicate(v, i)
				i++
			}
			yield(passed, nil)
		})
	}
}

func Find[T any](predicate func(T, int) bool) OperatorFunc[T, T] {
	return func(input Observable[T]) Observable[T] {
		return (ObservableFunc[T])(func(yield func(T, error) bool) {
			var i int
			for v, err := range input.Subscribe() {
				if err != nil {
					var zero T
					yield(zero, err)
					return
				}
				if predicate(v, i) {
					yield(v, nil)
					return
				}
				i++
			}
			var zero T
			yield(zero, ErrNotFound)
		})
	}
}

func FindIndex[T any](predicate func(T, int) bool) OperatorFunc[T, int] {
	return func(input Observable[T]) Observable[int] {
		return (ObservableFunc[int])(func(yield func(int, error) bool) {
			var i int
			for v, err := range input.Subscribe() {
				if err != nil {
					yield(-1, err)
					return
				} else {
					if predicate(v, i) {
						yield(i, nil)
						return
					}
				}
				i++
			}
			yield(-1, nil)
		})
	}
}

func IsEmpty[T any]() OperatorFunc[T, bool] {
	return func(input Observable[T]) Observable[bool] {
		return (ObservableFunc[bool])(func(yield func(bool, error) bool) {
			for _, err := range input.Subscribe() {
				if err != nil {
					yield(false, err)
					return
				}
				yield(false, nil)
				return
			}
			yield(true, nil)
		})
	}
}
