package rxgo

import (
	"time"
)

func DebounceTime[I, O any](d time.Duration) OperatorFunc[I, O] {
	return func(input Observable[I]) Observable[O] {
		return func(yield func(O, error) bool) {

		}
	}
}
