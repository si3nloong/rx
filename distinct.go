package rxgo

import (
	"iter"
	"reflect"
)

func DistinctUntilChanged[T any](comparator ...func(prev, curr T) bool) OperatorFunc[T, T] {
	return func(input Observable[T]) Observable[T] {
		return func(yield func(T, error) bool) {
			next, stop := iter.Pull2((iter.Seq2[T, error])(input))
			defer stop()

			latestValue, err, ok := next()
			if err != nil {
				var zero T
				yield(zero, err)
				return
			} else if !ok {
				yield(latestValue, nil)
				return
			} else {
				if !yield(latestValue, nil) {
					return
				}
			}

			fn := func(prev, curr T) bool {
				return reflect.DeepEqual(prev, curr)
			}
			if len(comparator) > 0 {
				fn = comparator[0]
			}

			for {
				v, err, ok := next()
				if err != nil {
					yield(v, err)
					return
				} else if !ok {
					return
				} else {
					if fn(latestValue, v) {
						continue
					}
					if !yield(v, nil) {
						return
					}
					latestValue = v
				}
			}
		}
	}
}
