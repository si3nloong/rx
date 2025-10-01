package rxgo

import (
	"iter"
	"time"
)

func Delay[T any](d time.Duration) OperatorFunc[T, T] {
	return func(input Observable[T]) Observable[T] {
		return func(yield func(T, error) bool) {
			t := time.NewTimer(d)
			defer t.Stop()

			<-t.C
			t.Stop()

			next, stop := iter.Pull2((iter.Seq2[T, error])(input))
			defer stop()

			for {
				v, err, ok := next()
				if err != nil {
					var zero T
					yield(zero, err)
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
}

func DebounceTime[I, O any](d time.Duration) OperatorFunc[I, O] {
	return func(input Observable[I]) Observable[O] {
		return func(yield func(O, error) bool) {

		}
	}
}
