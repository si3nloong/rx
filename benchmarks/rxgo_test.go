package rxgo

import (
	"context"
	"strconv"
	"testing"

	rxgo2 "github.com/reactivex/rxgo/v2"
	"github.com/samber/ro"
	"github.com/si3nloong/rx"
)

func BenchmarkRx(b *testing.B) {
	arr := make([]int, 0, 100)
	for i := 0; i < 100; i++ {
		arr = append(arr, i)
	}

	// b.Run("Go", func(b *testing.B) {
	// 	for i := 0; i < b.N; i++ {
	// 		newArr := make([]any, 0, len(arr))
	// 		for _, a := range arr {
	// 			if a%2 == 0 {
	// 				newArr = append(newArr, a)
	// 			}
	// 		}
	// 		for i, a := range newArr {
	// 			newArr[i] = "num" + strconv.Itoa(a.(int))
	// 		}
	// 		_ = newArr
	// 	}

	// 	b.Run("Filter", func(b *testing.B) {
	// 		for i := 0; i < b.N; i++ {
	// 			newArr := make([]any, 0, len(arr))
	// 			for _, a := range arr {
	// 				if a%2 == 0 {
	// 					newArr = append(newArr, a)
	// 				}
	// 			}
	// 			_ = newArr
	// 		}
	// 	})
	// })
	b.Run("Range", func(b *testing.B) {
		b.Run("RxGo", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				rxgo2.Range(1, 100).
					DoOnNext(func(i interface{}) {})
			}
		})
		b.Run("ro", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				ro.Range(1, 100).Subscribe(ro.OnNext(func(v int64) {}))
			}
		})
		b.Run("rxgo", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				rx.Range(1, 100).SubscribeOn(
					func(v int) {},
					func(err error) {},
					func() {},
				)
			}
		})
	})
	b.Run("ToSlice", func(b *testing.B) {
		b.Run("RxGo", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				rxgo2.Range(1, 100).
					ToSlice(100)
			}
		})
		b.Run("ro", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				ro.Pipe1(
					ro.Range(1, 100),
					ro.ToSlice[int64](),
				).Subscribe(ro.OnNext(func(v []int64) {}))
			}
		})
		b.Run("rxgo", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				rx.Pipe1(
					rx.Range(1, 100),
					rx.ToSlice[int](),
				).SubscribeOn(
					func(v []int) {},
					func(err error) {},
					func() {},
				)
			}
		})
	})
	b.Run("Filter", func(b *testing.B) {
		b.Run("RxGo", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				rxgo2.Just(arr)().
					Filter(func(v any) bool {
						return v.(int)%2 == 0
					}).
					DoOnNext(func(i interface{}) {})
			}
		})
		b.Run("ro", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				ro.Pipe1(
					ro.FromSlice(arr),
					ro.Filter(func(v int) bool {
						return v%2 == 0
					}),
				).Subscribe(ro.OnNext(func(v int) {}))
			}
		})
		b.Run("rxgo", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				rx.Pipe1(
					rx.From[int](arr),
					rx.Filter(func(v int) bool {
						return v%2 == 0
					}),
				).SubscribeOn(
					func(v int) {},
					func(err error) {},
					func() {},
				)
			}
		})
	})
	b.Run("IgnoreElements", func(b *testing.B) {
		b.Run("RxGo", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				rxgo2.Just(arr)().
					IgnoreElements().
					DoOnNext(func(i interface{}) {})
			}
		})
		b.Run("ro", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				ro.Pipe1(
					ro.FromSlice(arr),
					ro.IgnoreElements[int](),
				).Subscribe(ro.OnNext(func(i int) {}))
			}
		})
		b.Run("rxgo", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				rx.Pipe1(
					rx.From[int](arr),
					rx.IgnoreElements[int](),
				).SubscribeOn(
					func(v int) {},
					func(err error) {},
					func() {},
				)
			}
		})
	})
	b.Run("Skip", func(b *testing.B) {
		b.Run("RxGo", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				rxgo2.Just(arr)().
					Skip(50).
					DoOnNext(func(i interface{}) {})
			}
		})
		b.Run("ro", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				ro.Pipe1(
					ro.FromSlice(arr),
					ro.Skip[int](50),
				).Subscribe(ro.OnNext(func(v int) {}))
			}
		})
		b.Run("rxgo", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				rx.Pipe1(
					rx.From[int](arr),
					rx.Skip[int](50),
				).SubscribeOn(
					func(v int) {},
					func(err error) {},
					func() {},
				)
			}
		})
	})
	b.Run("Take", func(b *testing.B) {
		b.Run("RxGo", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				rxgo2.Just(arr)().
					Take(2).
					DoOnNext(func(i interface{}) {})
			}
		})
		b.Run("ro", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				ro.Pipe1(
					ro.FromSlice(arr),
					ro.Take[int](2),
				).Subscribe(ro.OnNext(func(i int) {}))
			}
		})
		b.Run("rxgo", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				rx.Pipe1(
					rx.From[int](arr),
					rx.Take[int](2),
				).SubscribeOn(
					func(v int) {},
					func(err error) {},
					func() {},
				)
			}
		})
	})
	b.Run("TakeLast", func(b *testing.B) {
		b.Run("RxGo", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				rxgo2.Just(arr)().
					TakeLast(1).
					DoOnNext(func(i interface{}) {})
			}
		})
		b.Run("ro", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				ro.Pipe1(
					ro.FromSlice(arr),
					ro.TakeLast[int](1),
				).Subscribe(ro.OnNext(func(i int) {}))
			}
		})
		b.Run("rxgo", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				rx.Pipe1(
					rx.From[int](arr),
					rx.TakeLast[int](1),
				).SubscribeOn(
					func(v int) {},
					func(err error) {},
					func() {},
				)
			}
		})
	})
	b.Run("Map", func(b *testing.B) {
		b.Run("RxGo", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				rxgo2.Just(arr)().
					Map(func(ctx context.Context, v any) (any, error) {
						return strconv.Itoa(v.(int)), nil
					}).
					DoOnNext(func(i interface{}) {})
			}
		})
		b.Run("ro", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				ro.Pipe1(
					ro.FromSlice(arr),
					ro.Map(func(v int) string {
						return strconv.Itoa(v)
					}),
				).Subscribe(ro.OnNext(func(v string) {}))
			}
		})
		b.Run("rxgo", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				rx.Pipe1(
					rx.From[int](arr),
					rx.Map(func(v int, _ int) string {
						return strconv.Itoa(v)
					}),
				).SubscribeOn(
					func(v string) {},
					func(err error) {},
					func() {},
				)
			}
		})
	})
	// b.Run("ro", func(b *testing.B) {
	// 	for i := 0; i < b.N; i++ {
	// 		ro.Pipe2(
	// 			ro.FromSlice(arr),
	// 			ro.Filter(func(v int) bool {
	// 				return v%2 == 0
	// 			}),
	// 			ro.Map(func(v int) string {
	// 				return "num" + strconv.Itoa(v)
	// 			}),
	// 		).Subscribe(ro.OnNext(func(v string) {}))
	// 	}

	// })
	b.Run("rxgo", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			rx.Pipe2(
				rx.From[int](arr),
				rx.Filter(func(v int) bool {
					return v%2 == 0
				}),
				rx.Map(func(v int, _ int) string {
					return "num" + strconv.Itoa(v)
				}),
			).SubscribeOn(
				func(v string) {},
				func(err error) {},
				func() {},
			)
		}
	})
}
