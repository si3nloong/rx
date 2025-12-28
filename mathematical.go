package rx

// Count counts the number of emissions on the source and emits that number when the source completes.
func Count[T Number](predicate ...func(value T, index int) bool) OperatorFunc[T, T] {
	return func(input Observable[T]) Observable[T] {
		return (ObservableFunc[T])(func(yield func(T, error) bool) {
			var count T
			if len(predicate) > 0 {
				var i int
				fn := predicate[0]
				for v, err := range input.Subscribe() {
					if err != nil {
						var zero T
						yield(zero, err)
						return
					}
					if fn(v, i) {
						count++
					}
					i++
				}
			} else {
				for _, err := range input.Subscribe() {
					if err != nil {
						var zero T
						yield(zero, err)
						return
					}
					count++
				}
			}

			yield(count, nil)
			count = 0
		})
	}
}

// Min emits the item from the source Observable that had the minimum value.
func Min[T Number]() OperatorFunc[T, T] {
	return func(input Observable[T]) Observable[T] {
		return (ObservableFunc[T])(func(yield func(T, error) bool) {
			var minValue T
			for v, err := range input.Subscribe() {
				if err != nil {
					var zero T
					yield(zero, err)
					return
				}
				minValue = min(minValue, v)
			}
			yield(minValue, nil)
		})
	}
}

// Max emits the item from the source Observable that had the maximum value.
func Max[T Number]() OperatorFunc[T, T] {
	return func(input Observable[T]) Observable[T] {
		return (ObservableFunc[T])(func(yield func(T, error) bool) {
			var maxValue T
			for v, err := range input.Subscribe() {
				if err != nil {
					var zero T
					yield(zero, err)
					return
				}
				maxValue = max(maxValue, v)
			}
			yield(maxValue, nil)
		})
	}
}

// Reduce applies an accumulator function over the source Observable, and returns the accumulated result when the source completes, given an optional seed value.
func Reduce[V, A any](accumulator func(acc A, value V, index int) A, seed A) OperatorFunc[V, A] {
	return func(input Observable[V]) Observable[A] {
		return (ObservableFunc[A])(func(yield func(A, error) bool) {
			var (
				acc = seed
				i   int
			)
			for v, err := range input.Subscribe() {
				if err != nil {
					var zero A
					yield(zero, err)
					return
				} else {
					acc = accumulator(acc, v, i)
					i++
				}
			}
			yield(acc, nil)
		})
	}
}
