package rxgo

import (
	"reflect"
	"testing"
)

func TestDistinct(t *testing.T) {
	Pipe2(
		Of([]int{1, 1, 1, 2, 2, 2, 1, 1, 3, 3}),
		DistinctUntilChanged[int](),
		ToSlice[int](),
	).Subscribe(func(v []int) {
		if !reflect.DeepEqual(v, []int{1, 2, 1, 3}) {
			t.Failed()
		}
	}, func(err error) {}, func() {
		t.Failed()
	})
}
