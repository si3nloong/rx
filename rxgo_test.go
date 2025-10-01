package rxgo

import (
	"strconv"
	"testing"
)

func BenchmarkRx(b *testing.B) {
	arr := []int{1, 2, 3, 4, 7, 8, 9, 10, 11, 14}

	b.Run("Rx", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			Pipe2(
				Of(arr),
				Filter(func(v int) bool {
					return v%2 == 0
				}),
				Map(func(v int, _ int) string {
					return "num" + strconv.Itoa(v)
				}),
			).Subscribe(
				func(v string) {},
				func(err error) {},
				func() {},
			)
		}
	})

	b.Run("Native", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			newArr := make([]any, 0, len(arr))
			for _, a := range arr {
				if a%2 == 0 {
					newArr = append(newArr, a)
				}
			}
			for i, a := range newArr {
				newArr[i] = "num" + strconv.Itoa(a.(int))
			}
			_ = newArr
		}
	})
}
