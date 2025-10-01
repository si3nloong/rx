package rxgo

type OperatorFunc[I any, O any] func(Observable[I]) Observable[O]
