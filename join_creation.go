package rxgo

import (
	"context"
	"iter"
	"sync"
	"sync/atomic"
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
					next, stop := iter.Pull2(input.Subscribe())
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
			next, stop := iter.Pull2(inputs[0].Subscribe())

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
					next, stop := iter.Pull2(input.Subscribe())
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
					next, stop := iter.Pull2(input.Subscribe())
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
		var wg sync.WaitGroup
		var once sync.Once
		var selectedIndex int32

		for i := range inputs {
			wg.Go(func(index int, input Observable[T]) func() {
				return func() {
					next, stop := iter.Pull2(input.Subscribe())
					defer stop()

					// Peek the first emission
					v, err, ok := next()
					if !ok {
						return
					} else {
						once.Do(func() {
							atomic.StoreInt32(&selectedIndex, (int32)(index))
						})
						if atomic.LoadInt32(&selectedIndex) != (int32)(index) {
							return
						}
						// If one of the used source observable throws an errors before a first notification the race operator will also throw an error, no matter if another source observable could potentially win the race.
						if err != nil {
							var zero T
							yield(zero, err)
							return
						}

						// As soon as one of the source observables emits a value, the result unsubscribes from the other sources.
						// The resulting observable will forward all notifications, including error and completion, from the "winning" source observable.
						if !yield(v, nil) {
							return
						}

						for {
							v, err, ok := next()
							if err != nil {
								return
							} else if !ok {
								return
							} else {
								if !yield(v, nil) {
									return
								}
							}
						}
					}
				}
			}(i, inputs[i]))
		}
		wg.Wait()
	})
}
