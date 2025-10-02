package rxgo

import (
	"context"
	"iter"
	"sync"
)

func Buffer[T any, I any](closingNotifier Observable[I]) OperatorFunc[T, []T] {
	return func(input Observable[T]) Observable[[]T] {
		return func(yield func([]T, error) bool) {
			var (
				buffer = make([]T, 0)
				rw     sync.RWMutex
			)

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			go func() {
				next, stop := iter.Pull2((iter.Seq2[T, error])(input))
				defer stop()

				for {
					select {
					case <-ctx.Done():
						return
					default:
						v, err, ok := next()
						if err != nil {
							yield(nil, err)
							return
						} else if !ok {
							return
						} else {
							rw.Lock()
							buffer = append(buffer, v)
							rw.Unlock()
						}
					}
				}
			}()

			next, stop := iter.Pull2((iter.Seq2[I, error])(closingNotifier))
			defer stop()

			for {
				if _, err, ok := next(); err != nil {
					cancel()
					yield(nil, err)
					return
				} else if !ok {
					cancel()
					return
				} else {
					rw.Lock()
					if !yield(buffer, nil) {
						rw.Unlock()
						return
					}
					buffer = make([]T, 0)
					rw.Unlock()
				}
			}
		}
	}
}

func Map[I, O any](fn func(v I, index int) O) OperatorFunc[I, O] {
	return func(input Observable[I]) Observable[O] {
		return func(yield func(O, error) bool) {
			next, stop := iter.Pull2((iter.Seq2[I, error])(input))
			defer stop()

			var i int
			for {
				v, err, ok := next()
				if err != nil {
					var o O
					yield(o, err)
					return
				} else if !ok {
					return
				} else {
					if !yield(fn(v, i), nil) {
						return
					}
				}
				i++
			}
		}
	}
}

func Map2[I, O any](fn func(v I, index int) (O, error)) OperatorFunc[I, O] {
	return func(input Observable[I]) Observable[O] {
		return func(yield func(O, error) bool) {
			next, stop := iter.Pull2((iter.Seq2[I, error])(input))
			defer stop()

			var i int
			for {
				v, err, ok := next()
				if err != nil {
					var zero O
					yield(zero, err)
					return
				} else if !ok {
					return
				} else {
					o, err := fn(v, i)
					if err != nil {
						yield(o, err)
						return
					}
					if !yield(o, nil) {
						return
					}
				}
				i++
			}
		}
	}
}

func ConcatMap[I, O any](project func(v I, index int) Observable[O]) OperatorFunc[I, O] {
	return func(input Observable[I]) Observable[O] {
		return func(yield func(O, error) bool) {
			next, stop := iter.Pull2((iter.Seq2[I, error])(input))
			defer stop()

			var i int
			for {
				v, err, ok := next()
				if err != nil {
					var zero O
					yield(zero, err)
					return
				} else if !ok {
					return
				} else {
					next2, stop2 := iter.Pull2((iter.Seq2[O, error])(project(v, i)))
				loop2:
					for {
						v2, err2, ok2 := next2()
						if err2 != nil {
							stop2()
							var zero O
							yield(zero, err2)
							return
						} else if !ok2 {
							break loop2
						} else {
							if !yield(v2, nil) {
								stop2()
								return
							}
						}
					}
					stop2()
				}
				i++
			}
		}
	}
}

func SwitchMap[I, O any](fn func(v I, index int) Observable[O]) OperatorFunc[I, O] {
	return func(input Observable[I]) Observable[O] {
		return func(yield func(O, error) bool) {
			next, stop := iter.Pull2((iter.Seq2[I, error])(input))
			defer stop()

			var i int

			for {
				v, err, ok := next()
				if err != nil {
					var zero O
					yield(zero, err)
					return
				} else if !ok {
					return
				} else {
					next2, stop2 := iter.Pull2((iter.Seq2[O, error])(fn(v, i)))

				loop2:
					for {
						v2, err2, ok2 := next2()
						if err2 != nil {
							var zero O
							yield(zero, err2)
							return
						} else if !ok2 {
							break loop2
						} else {
							if !yield(v2, nil) {
								return
							}
						}
					}
					stop2()
				}
				i++
			}
		}
	}
}

func MergeMap[T any](fn func(v T, index int) Observable[T]) OperatorFunc[T, T] {
	return func(input Observable[T]) Observable[T] {
		return func(yield func(T, error) bool) {
			next, stop := iter.Pull2((iter.Seq2[T, error])(input))
			defer stop()

			var i int

			for {
				v, err, ok := next()
				if err != nil {
					var zero T
					yield(zero, err)
					return
				} else if !ok {
					return
				} else {
					next2, stop2 := iter.Pull2((iter.Seq2[T, error])(fn(v, i)))

				loop2:
					for {
						v2, err2, ok2 := next2()
						if err2 != nil {
							var zero T
							yield(zero, err2)
							return
						} else if !ok2 {
							break loop2
						} else {
							if !yield(v2, nil) {
								return
							}
						}
					}
					stop2()
				}
				i++
			}
		}
	}
}

func Pairwise[T any]() OperatorFunc[T, [2]T] {
	return func(input Observable[T]) Observable[[2]T] {
		return func(yield func([2]T, error) bool) {
			next, stop := iter.Pull2((iter.Seq2[T, error])(input))
			defer stop()

			var n int
			var pair [2]T

			for {
				v, err, ok := next()
				if err != nil {
					var zero [2]T
					yield(zero, err)
					return
				} else if !ok {
					return
				} else {
					if n%2 == 0 {
						pair[1] = v
						if !yield(pair, nil) {
							return
						}
						pair = [2]T{}
					} else {
						pair[0] = v
					}
					n++
				}
			}
		}
	}
}

func Scan[V, A any](accumulator func(acc A, value V, index int) A, seed A) OperatorFunc[V, A] {
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
					if !yield(acc, nil) {
						return
					}
					i++
				}
			}
		}
	}
}
