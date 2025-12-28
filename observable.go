package rx

import (
	"iter"
)

// ObservableFunc is a function type that implements the Observable interface.
// It allows casting a function matching iter.Seq2 signature to an Observable.
type ObservableFunc[T any] iter.Seq2[T, error]

// All returns an iterator over the elements of the Observable.
// It yields values until an error occurs or the Observable completes.
// If an error occurs, it yields the zero value of T and stops.
func (fn ObservableFunc[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		for v, err := range (iter.Seq2[T, error])(fn) {
			if err != nil {
				var zero T
				yield(zero)
				return
			} else {
				if !yield(v) {
					return
				}
			}
		}
	}
}

// Subscribe returns an iterator that yields values and errors from the Observable.
// It delegates to the underlying iter.Seq2 function.
func (fn ObservableFunc[T]) Subscribe() iter.Seq2[T, error] {
	return (iter.Seq2[T, error])(fn)
}

// SubscribeOn subscribes to the Observable and executes the provided callbacks for each event.
// onNext is called for each value emitted.
// onError is called if an error occurs.
// onComplete is called when the Observable completes successfully.
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

// Pipe1 pipes the input Observable through 1 operator function.
//
// Example:
//
//	rx.Pipe1(
//		rx.Of(1, 2, 3),
//		rx.Map(func(i int, _ int) int { return i * 2 }),
//	)
func Pipe1[I, O1 any](
	input Observable[I],
	f1 OperatorFunc[I, O1],
) Observable[O1] {
	return Observable[O1](f1(input))
}

// Pipe2 pipes the input Observable through 2 operator functions.
//
// Example:
//
//	rx.Pipe2(
//		rx.Of(1, 2, 3),
//		rx.Map(func(i int, _ int) int { return i * 2 }),
//		rx.Filter(func(i int) bool { return i > 2 }),
//	)
func Pipe2[I, O1, O2 any](
	input Observable[I],
	f1 OperatorFunc[I, O1],
	f2 OperatorFunc[O1, O2],
) Observable[O2] {
	return Observable[O2](f2(f1(input)))
}

// Pipe3 pipes the input Observable through 3 operator functions.
func Pipe3[I, O1, O2, O3 any](
	input Observable[I],
	f1 OperatorFunc[I, O1],
	f2 OperatorFunc[O1, O2],
	f3 OperatorFunc[O2, O3],
) Observable[O3] {
	return Observable[O3](f3(f2(f1(input))))
}

// Pipe4 pipes the input Observable through 4 operator functions.
func Pipe4[I, O1, O2, O3, O4 any](
	input Observable[I],
	f1 OperatorFunc[I, O1],
	f2 OperatorFunc[O1, O2],
	f3 OperatorFunc[O2, O3],
	f4 OperatorFunc[O3, O4],
) Observable[O4] {
	return Observable[O4](f4(f3(f2(f1(input)))))
}

// Pipe5 pipes the input Observable through 5 operator functions.
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

// Pipe6 pipes the input Observable through 6 operator functions.
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

// Pipe7 pipes the input Observable through 7 operator functions.
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

// Pipe8 pipes the input Observable through 8 operator functions.
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

// Pipe9 pipes the input Observable through 9 operator functions.
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

// Pipe10 pipes the input Observable through 10 operator functions.
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
