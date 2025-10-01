package rxgo

import "iter"

func ThrottleTime[I, O any](fn func(v I, index int) O) OperatorFunc[I, O] {
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
					if !yield(fn(v, i), nil) {
						return
					}
				}
				i++
			}
		}
	}
}
