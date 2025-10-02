package rxgo

import (
	"iter"
)

func Range[T Number](start, count T) Observable[T] {
	return func(yield func(T, error) bool) {
		for ; start < count; start++ {
			if !yield(start, nil) {
				return
			}
		}
	}
}

func Count[T Number](predicate ...func(value T, index int) bool) OperatorFunc[T, T] {
	return func(input Observable[T]) Observable[T] {
		return func(yield func(T, error) bool) {
			next, stop := iter.Pull2((iter.Seq2[T, error])(input))
			defer stop()

			var count T
			if len(predicate) > 0 {
				var i int
				fn := predicate[0]

				for {
					v, err, ok := next()
					if err != nil {
						var zero T
						yield(zero, err)
						return
					} else if !ok {
						break
					} else {
						if fn(v, i) {
							count++
						}
						i++
					}
				}
			} else {
				for {
					if _, err, ok := next(); err != nil {
						var zero T
						yield(zero, err)
						return
					} else if !ok {
						break
					} else {
						count++
					}
				}
			}

			yield(count, nil)
		}
	}
}

func Min[T Number]() OperatorFunc[T, T] {
	return func(input Observable[T]) Observable[T] {
		return func(yield func(T, error) bool) {
			next, stop := iter.Pull2((iter.Seq2[T, error])(input))
			defer stop()

			minValue, err, ok := next()
			if err != nil {
				var zero T
				yield(zero, err)
				return
			} else if !ok {
				yield(minValue, nil)
				return
			}

			for {
				v, err, ok := next()
				if err != nil {
					var zero T
					yield(zero, err)
					return
				} else if !ok {
					break
				} else {
					minValue = min(minValue, v)
				}
			}

			if !yield(minValue, nil) {
				return
			}
		}
	}
}

func Max[T Number]() OperatorFunc[T, T] {
	return func(input Observable[T]) Observable[T] {
		return func(yield func(T, error) bool) {
			next, stop := iter.Pull2((iter.Seq2[T, error])(input))
			defer stop()

			maxValue, err, ok := next()
			if err != nil {
				var zero T
				yield(zero, err)
				return
			} else if !ok {
				yield(maxValue, nil)
				return
			}

			for {
				v, err, ok := next()
				if err != nil {
					var zero T
					yield(zero, err)
					return
				} else if !ok {
					break
				} else {
					maxValue = max(maxValue, v)
				}
			}

			if !yield(maxValue, nil) {
				return
			}
		}
	}
}

func Reduce[V, A any](accumulator func(acc A, value V, index int) A, seed A) OperatorFunc[V, A] {
	return func(input Observable[V]) Observable[A] {
		return func(yield func(A, error) bool) {
			next, stop := iter.Pull2((iter.Seq2[V, error])(input))
			defer stop()

			var (
				acc = seed
				i   int
			)

			for {
				v, err, ok := next()
				if err != nil {
					var zero A
					yield(zero, err)
					return
				} else if !ok {
					break
				} else {
					acc = accumulator(acc, v, i)
					i++
				}
			}

			yield(acc, nil)
		}
	}
}
