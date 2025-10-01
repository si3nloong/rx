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
