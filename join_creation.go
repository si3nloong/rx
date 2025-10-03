package rxgo

import (
	"context"
	"iter"
	"sync"
)

func CombineLatest[T any](inputs ...Observable[T]) Observable[[]T] {
	if len(inputs) < 2 {
		panic(`CombineLatest required at least 2 observable`)
	}
	return (ObservableFunc[[]T])(func(yield func([]T, error) bool) {
		inputCount := len(inputs)
		ch := make(chan state[T], 1)
		defer close(ch)

		var wg sync.WaitGroup
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		for i := range inputs {
			wg.Go(func(index int, input Observable[T]) func() {
				return func() {
					next, stop := iter.Pull2((iter.Seq2[T, error])(input.Subscribe()))
					defer stop()

					for {
						v, err, ok := next()
						select {
						case <-ctx.Done():
							return
						case ch <- state[T]{index, v, err, ok}:
							if err != nil || !ok {
								return
							}
						}
					}
				}
			}(i, inputs[i]))
		}

		var counter int
		idxCache := make(map[int]struct{})
		defer clear(idxCache)
		results := make([]T, len(inputs))

	loop:
		for o := range ch {
			if o.err != nil {
				yield(nil, o.err)
				return
			} else if !o.ok {
				counter++
				if counter >= inputCount {
					break loop
				}
			} else {
				idxCache[o.idx] = struct{}{}
				results[o.idx] = o.v
				if len(idxCache) >= inputCount {
					if !yield(results, nil) {
						return
					}
				}
			}
		}

		wg.Wait()
	})
}

func Concat[T any](inputs ...Observable[T]) Observable[T] {
	if len(inputs) < 2 {
		panic(`Concat required at least 2 observable`)
	}
	return (ObservableFunc[T])(func(yield func(T, error) bool) {
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

func Merge[T any](inputs ...Observable[T]) Observable[T] {
	if len(inputs) < 2 {
		panic(`Merge required at least 2 observable`)
	}
	return (ObservableFunc[T])(func(yield func(T, error) bool) {
		ch := make(chan T, 1)

		g, ctx := WithContext(context.Background())
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
							ch <- v
						}
					}
				}
			}(i, v))
		}
		go func() {
			defer close(ch)

			for {
				select {
				case <-ctx.Done():
					return
				case v := <-ch:
					if !yield(v, nil) {
						return
					}
				}
			}
		}()
		if err := g.Wait(); err != nil {
			var zero T
			yield(zero, err)
			return
		}
	})
}

func Race[T any](inputs ...Observable[T]) Observable[T] {
	if len(inputs) < 2 {
		panic(`Race required at least 2 observable`)
	}
	return (ObservableFunc[T])(func(yield func(T, error) bool) {
	})
}

func Zip[T any](inputs ...Observable[T]) Observable[[]T] {
	if len(inputs) < 2 {
		panic(`Zip required at least 2 observable`)
	}
	return (ObservableFunc[[]T])(func(yield func([]T, error) bool) {

	})
}
