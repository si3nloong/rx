package rx

import (
	"reflect"
	"testing"
)

func TestFiltering(t *testing.T) {
	t.Run("DistinctUntilChanged", func(t *testing.T) {
		Pipe2(
			Of(1, 1, 1, 2, 2, 2, 1, 1, 3, 3),
			DistinctUntilChanged[int](),
			ToSlice[int](),
		).SubscribeOn(func(v []int) {
			if !reflect.DeepEqual(v, []int{1, 2, 1, 3}) {
				t.Failed()
			}
		}, func(err error) {}, func() {
			t.Failed()
		})
	})
}
