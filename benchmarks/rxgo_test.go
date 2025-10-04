package rxgo

import (
	"context"
	"strconv"
	"testing"

	rxgo2 "github.com/reactivex/rxgo/v2"
	"github.com/si3nloong/rxgo"
)

func BenchmarkRx(b *testing.B) {
	arr := make([]int, 0, 100)
	for i := 0; i < 100; i++ {
		arr = append(arr, i)
	}

	b.Run("Go", func(b *testing.B) {
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

		b.Run("Filter", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				newArr := make([]any, 0, len(arr))
				for _, a := range arr {
					if a%2 == 0 {
						newArr = append(newArr, a)
					}
				}
				_ = newArr
			}
		})
	})
	b.Run("rxgo", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			rxgo.Pipe2(
				rxgo.From(arr),
				rxgo.Filter(func(v int) bool {
					return v%2 == 0
				}),
				rxgo.Map(func(v int, _ int) string {
					return "num" + strconv.Itoa(v)
				}),
			).SubscribeOn(
				func(v string) {},
				func(err error) {},
				func() {},
			)
		}

		b.Run("Range", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				rxgo.Range(1, 100).SubscribeOn(
					func(v int) {},
					func(err error) {},
					func() {},
				)
			}
		})
		b.Run("ToSlice", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				rxgo.Pipe1(
					rxgo.Range(1, 100),
					rxgo.ToSlice[int](),
				).SubscribeOn(
					func(v []int) {},
					func(err error) {},
					func() {},
				)
			}
		})
		b.Run("Skip", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				rxgo.Pipe1(
					rxgo.From(arr),
					rxgo.Skip[int](50),
				).SubscribeOn(
					func(v int) {},
					func(err error) {},
					func() {},
				)
			}
		})
		b.Run("Filter", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				rxgo.Pipe1(
					rxgo.From(arr),
					rxgo.Filter(func(v int) bool {
						return v%2 == 0
					}),
				).SubscribeOn(
					func(v int) {},
					func(err error) {},
					func() {},
				)
			}
		})
		b.Run("Take", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				rxgo.Pipe1(
					rxgo.From(arr),
					rxgo.Take[int](2),
				).SubscribeOn(
					func(v int) {},
					func(err error) {},
					func() {},
				)
			}
		})
		b.Run("TakeLast", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				rxgo.Pipe1(
					rxgo.From(arr),
					rxgo.TakeLast[int](1),
				).SubscribeOn(
					func(v int) {},
					func(err error) {},
					func() {},
				)
			}
		})
		b.Run("Map", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				rxgo.Pipe1(
					rxgo.From(arr),
					rxgo.Map(func(v int, _ int) string {
						return strconv.Itoa(v)
					}),
				).SubscribeOn(
					func(v string) {},
					func(err error) {},
					func() {},
				)
			}
		})
		b.Run("IgnoreElements", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				rxgo.Pipe1(
					rxgo.From(arr),
					rxgo.IgnoreElements[int](),
				).SubscribeOn(
					func(v int) {},
					func(err error) {},
					func() {},
				)
			}
		})
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
				DoOnNext(func(i interface{}) {})
		}

		b.Run("Range", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				rxgo2.Range(1, 100).
					DoOnNext(func(i interface{}) {})
			}
		})
		b.Run("ToSlice", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				rxgo2.Range(1, 100).
					ToSlice(100)
			}
		})
		b.Run("Skip", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				rxgo2.Just(arr)().
					Skip(50).
					DoOnNext(func(i interface{}) {})
			}
		})
		b.Run("Filter", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				rxgo2.Just(arr)().
					Filter(func(v any) bool {
						return v.(int)%2 == 0
					}).
					DoOnNext(func(i interface{}) {})
			}
		})
		b.Run("Take", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				rxgo2.Just(arr)().
					Take(2).
					DoOnNext(func(i interface{}) {})
			}
		})
		b.Run("TakeLast", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				rxgo2.Just(arr)().
					TakeLast(1).
					DoOnNext(func(i interface{}) {})
			}
		})
		b.Run("Map", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				rxgo2.Just(arr)().
					Map(func(ctx context.Context, v any) (any, error) {
						return strconv.Itoa(v.(int)), nil
					}).
					DoOnNext(func(i interface{}) {})
			}
		})
		b.Run("IgnoreElements", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				rxgo2.Just(arr)().
					IgnoreElements().
					DoOnNext(func(i interface{}) {})
			}
		})
	})
}
