package rxgo

import (
	"context"
	"strconv"
	"testing"

	rxgo2 "github.com/reactivex/rxgo/v2"
	"github.com/si3nloong/rxgo"
)

func BenchmarkRx(b *testing.B) {
	arr := []int{1, 2, 3, 4, 7, 8, 9, 10, 11, 14}

	b.Run("Go (native)", func(b *testing.B) {
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

	b.Run("rxgo", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			rxgo.Pipe2(
				rxgo.Of(arr),
				rxgo.Filter(func(v int) bool {
					return v%2 == 0
				}),
				rxgo.Map(func(v int, _ int) string {
					return "num" + strconv.Itoa(v)
				}),
			).Subscribe(
				func(v string) {},
				func(err error) {},
				func() {},
			)
		}
	})

	b.Run("RxGo", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			rxgo2.Just(arr)().
				Filter(func(v any) bool {
					return v.(int)%2 == 0
				}).
				Map(func(ctx context.Context, i interface{}) (interface{}, error) {
					return "num" + strconv.Itoa(i.(int)), nil
				}).
				DoOnNext(func(i interface{}) {
					// println(i)
				})
		}
	})
}
