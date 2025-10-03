package rxgo

import (
	"iter"
)

type ObservableFunc[T any] iter.Seq2[T, error]

func (fn ObservableFunc[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		next, stop := iter.Pull2((iter.Seq2[T, error])(fn))
		defer stop()

		for {
			v, err, ok := next()
			if err != nil {
				var zero T
				yield(zero)
				return
			} else if !ok {
				return
			} else {
				if !yield(v) {
					return
				}
			}
		}
	}
}

func (fn ObservableFunc[T]) Subscribe() iter.Seq2[T, error] {
	return (iter.Seq2[T, error])(fn)
}

func (fn ObservableFunc[T]) SubscribeOn(onNext func(v T), onError func(err error), onComplete func()) {
	next, stop := iter.Pull2((iter.Seq2[T, error])(fn))
	defer stop()

	for {
		if v, err, ok := next(); err != nil {
			onError(err)
			return
		} else if !ok {
			onComplete()
			return
		} else {
			onNext(v)
		}
	}
}

func Pipe1[I, O1 any](
	input Observable[I],
	f1 OperatorFunc[I, O1],
) Observable[O1] {
	return Observable[O1](f1(input))
}

func Pipe2[I, O1, O2 any](
	input Observable[I],
	f1 OperatorFunc[I, O1],
	f2 OperatorFunc[O1, O2],
) Observable[O2] {
	return Observable[O2](f2(f1(input)))
}

func Pipe3[I, O1, O2, O3 any](
	input Observable[I],
	f1 OperatorFunc[I, O1],
	f2 OperatorFunc[O1, O2],
	f3 OperatorFunc[O2, O3],
) Observable[O3] {
	return Observable[O3](f3(f2(f1(input))))
}

func Pipe4[I, O1, O2, O3, O4 any](
	input Observable[I],
	f1 OperatorFunc[I, O1],
	f2 OperatorFunc[O1, O2],
	f3 OperatorFunc[O2, O3],
	f4 OperatorFunc[O3, O4],
) Observable[O4] {
	return Observable[O4](f4(f3(f2(f1(input)))))
}

func Pipe5[I, O1, O2, O3, O4, O5 any](
	input Observable[I],
	f1 OperatorFunc[I, O1],
	f2 OperatorFunc[O1, O2],
	f3 OperatorFunc[O2, O3],
	f4 OperatorFunc[O3, O4],
	f5 OperatorFunc[O4, O5],
) Observable[O5] {
	return Observable[O5](f5(f4(f3(f2(f1(input))))))
}

func Pipe6[I, O1, O2, O3, O4, O5, O6 any](
	input Observable[I],
	f1 OperatorFunc[I, O1],
	f2 OperatorFunc[O1, O2],
	f3 OperatorFunc[O2, O3],
	f4 OperatorFunc[O3, O4],
	f5 OperatorFunc[O4, O5],
	f6 OperatorFunc[O5, O6],
) Observable[O6] {
	return Observable[O6](f6(f5(f4(f3(f2(f1(input)))))))
}

func Pipe7[I, O1, O2, O3, O4, O5, O6, O7 any](
	input Observable[I],
	f1 OperatorFunc[I, O1],
	f2 OperatorFunc[O1, O2],
	f3 OperatorFunc[O2, O3],
	f4 OperatorFunc[O3, O4],
	f5 OperatorFunc[O4, O5],
	f6 OperatorFunc[O5, O6],
	f7 OperatorFunc[O6, O7],
) Observable[O7] {
	return Observable[O7](f7(f6(f5(f4(f3(f2(f1(input))))))))
}

func Pipe8[I, O1, O2, O3, O4, O5, O6, O7, O8 any](
	input Observable[I],
	f1 OperatorFunc[I, O1],
	f2 OperatorFunc[O1, O2],
	f3 OperatorFunc[O2, O3],
	f4 OperatorFunc[O3, O4],
	f5 OperatorFunc[O4, O5],
	f6 OperatorFunc[O5, O6],
	f7 OperatorFunc[O6, O7],
	f8 OperatorFunc[O7, O8],
) Observable[O8] {
	return Observable[O8](f8(f7(f6(f5(f4(f3(f2(f1(input)))))))))
}

func Pipe9[I, O1, O2, O3, O4, O5, O6, O7, O8, O9 any](
	input Observable[I],
	f1 OperatorFunc[I, O1],
	f2 OperatorFunc[O1, O2],
	f3 OperatorFunc[O2, O3],
	f4 OperatorFunc[O3, O4],
	f5 OperatorFunc[O4, O5],
	f6 OperatorFunc[O5, O6],
	f7 OperatorFunc[O6, O7],
	f8 OperatorFunc[O7, O8],
	f9 OperatorFunc[O8, O9],
) Observable[O9] {
	return Observable[O9](f9(f8(f7(f6(f5(f4(f3(f2(f1(input))))))))))
}

func Pipe10[I, O1, O2, O3, O4, O5, O6, O7, O8, O9, O10 any](
	input Observable[I],
	f1 OperatorFunc[I, O1],
	f2 OperatorFunc[O1, O2],
	f3 OperatorFunc[O2, O3],
	f4 OperatorFunc[O3, O4],
	f5 OperatorFunc[O4, O5],
	f6 OperatorFunc[O5, O6],
	f7 OperatorFunc[O6, O7],
	f8 OperatorFunc[O7, O8],
	f9 OperatorFunc[O8, O9],
	f10 OperatorFunc[O9, O10],
) Observable[O10] {
	return Observable[O10](f10(f9(f8(f7(f6(f5(f4(f3(f2(f1(input)))))))))))
}
