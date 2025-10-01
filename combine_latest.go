package rxgo

import "iter"

type Tuple[A, B any] struct {
	a A
	b B
}

func (t Tuple[A, B]) Left() A {
	return t.a
}

func (t Tuple[A, B]) Right() B {
	return t.b
}

func CombineLatestWith[T1, T2 any](otherSource Observable[T2]) OperatorFunc[T1, Tuple[T1, T2]] {
	return func(input Observable[T1]) Observable[Tuple[T1, T2]] {
		return func(yield func(Tuple[T1, T2], error) bool) {
			next, stop := iter.Pull2((iter.Seq2[T1, error])(input))
			defer stop()

			for {
				v, err, ok := next()
				if err != nil {
					yield(Tuple[T1, T2]{}, err)
					return
				} else if !ok {
					return
				} else {
					if !yield(Tuple[T1, T2]{a: v}, nil) {
						return
					}
				}
			}
		}
	}
}
