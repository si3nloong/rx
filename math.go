package rxgo

import (
	"iter"

	"golang.org/x/exp/constraints"
)

func Range[T constraints.Integer](start, count T) Observable[T] {
	return func(yield func(T, error) bool) {
		for ; start < count; start++ {
			if !yield(start, nil) {
				return
			}
		}
	}
}

func Min[T constraints.Ordered]() OperatorFunc[T, T] {
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

func Max[T constraints.Ordered]() OperatorFunc[T, T] {
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
