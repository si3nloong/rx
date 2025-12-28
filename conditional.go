package rx

// DefaultIfEmpty emits a default value if the source Observable completes without emitting any value.
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

// Every returns an Observable that emits true or false to the Observer.
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

// Find emits only the first value emitted by the source Observable that meets some condition.
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

// FindIndex emits only the index of the first value emitted by the source Observable that meets some condition.
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

// IsEmpty returns an Observable that emits true if the source Observable is empty, otherwise false.
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
