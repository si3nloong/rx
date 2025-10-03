package rxgo

import (
	"iter"
)

func Concat[T any](inputs ...Observable[T]) Observable[T] {
	if len(inputs) < 2 {
		panic(`Concat required at least 2 observable`)
	}
	return ObservableFunc[T](func(yield func(T, error) bool) {
		for len(inputs) > 0 {
			next, stop := iter.Pull2((iter.Seq2[T, error])(inputs[0].Subscribe()))

		innerLoop:
			for {
				v, err, ok := next()
				if err != nil {
					stop()
					var zero T
					yield(zero, err)
					return
				} else if !ok {
					stop()
					break innerLoop
				} else {
					if !yield(v, nil) {
						stop()
						return
					}
				}
			}
			// Unshift the slice
			inputs = inputs[1:]
		}
	})
}

func ForkJoin[T any](inputs ...Observable[T]) Observable[[]T] {
	if len(inputs) < 2 {
		panic(`ForkJoin required at least 2 observable`)
	}
	return (ObservableFunc[[]T])(func(yield func([]T, error) bool) {
		g := Group{}
		results := make([]T, len(inputs))
		for i, v := range inputs {
			g.Go(func(index int, input Observable[T]) func() error {
				return func() error {
					next, stop := iter.Pull2((iter.Seq2[T, error])(input.Subscribe()))
					defer stop()

					for {
						v, err, ok := next()
						if err != nil {
							return err
						} else if !ok {
							return nil
						} else {
							results[index] = v
						}
					}
				}
			}(i, v))
		}

		if err := g.Wait(); err != nil {
			yield(nil, err)
			return
		} else {
			yield(results, nil)
		}
	})
}

func Zip[T any](inputs ...Observable[T]) Observable[T] {
	if len(inputs) < 2 {
		panic(`Zip required at least 2 observable`)
	}
	return (ObservableFunc[T])(func(yield func(T, error) bool) {
	})
}
