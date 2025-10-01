package rxgo

func Repeat[T any](count uint) OperatorFunc[T, T] {
	return func(input Observable[T]) Observable[T] {
		return func(yield func(T, error) bool) {

		}
	}
}
