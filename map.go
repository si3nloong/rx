package rxgo

import (
	"iter"
)

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
