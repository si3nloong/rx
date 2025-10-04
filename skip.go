package rxgo

import (
	"iter"
)

func Skip[T any](count uint) OperatorFunc[T, T] {
	return func(input Observable[T]) Observable[T] {
		return (ObservableFunc[T])(func(yield func(T, error) bool) {
			var skipCount uint
			for v, err := range input.Subscribe() {
				if err != nil {
					var zero T
					yield(zero, err)
					return
				} else {
					skipCount++
					if skipCount > count {
						if !yield(v, nil) {
							return
						}
					}
				}
			}
		})
	}
}

func SkipLast[T any](count uint) OperatorFunc[T, T] {
	return func(input Observable[T]) Observable[T] {
		return (ObservableFunc[T])(func(yield func(T, error) bool) {
			result := make([]T, 0, count)
			for v, err := range input.Subscribe() {
				if err != nil {
					var zero T
					yield(zero, err)
					return
				} else {
					result = append(result, v)
				}
			}

			if (uint)(len(result)) > count {
				result = result[:count]
				for _, v := range result {
					if !yield(v, nil) {
						return
					}
				}
			}
		})
	}
}

func SkipWhile[T any](fn func(T, int) bool) OperatorFunc[T, T] {
	return func(input Observable[T]) Observable[T] {
		return (ObservableFunc[T])(func(yield func(T, error) bool) {
			var i int
			for v, err := range input.Subscribe() {
				if err != nil {
					yield(v, err)
					return
				} else {
					if !fn(v, i) {
						if !yield(v, nil) {
							return
						}
					}
				}
				i++
			}
		})
	}
}

func SkipUntil[T, U any](notifier Observable[U]) OperatorFunc[T, T] {
	return func(input Observable[T]) Observable[T] {
		return (ObservableFunc[T])(func(yield func(T, error) bool) {
			ch := make(chan state[U], 1)
			go func() {
				next, stop := iter.Pull2(notifier.Subscribe())
				defer stop()

				v, err, ok := next()
				ch <- state[U]{0, v, err, ok}
			}()

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
