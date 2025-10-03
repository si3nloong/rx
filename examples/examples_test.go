package main_test

import (
	"errors"
	"testing"
	"time"

	"github.com/si3nloong/rxgo"
	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

func TestBuffer(t *testing.T) {
	defer goleak.VerifyNone(t)

	assertItems(t, rxgo.Pipe1(
		rxgo.Of([]int{1, 3, 4, 5, 9}),
		rxgo.BufferCount[int](2),
	), [][]int{
		{1, 3},
		{4, 5},
		{9},
	})
}

func TestCombineLatest(t *testing.T) {
	defer goleak.VerifyNone(t)

	assertItems(t, rxgo.CombineLatest(
		rxgo.Pipe2(rxgo.Of([]int{1}), rxgo.Delay[int](time.Second), rxgo.StartWith(0)),
		rxgo.Pipe2(rxgo.Of([]int{5}), rxgo.Delay[int](time.Second*5), rxgo.StartWith(0)),
		rxgo.Pipe2(rxgo.Of([]int{10}), rxgo.Delay[int](time.Second*10), rxgo.StartWith(0)),
	), [][]int{
		{0, 0, 0},
		{1, 0, 0},
		{1, 5, 0},
		{1, 5, 10},
	})
}

func TestCatchError(t *testing.T) {
	defer goleak.VerifyNone(t)

	assertItem(t, rxgo.Pipe3(
		rxgo.Range(1, 20),
		rxgo.Filter(func(v int) bool {
			return v%2 == 0
		}),
		rxgo.Map2(func(v int, _ int) (int, error) {
			if v > 10 {
				return 0, errors.New(`stop la`)
			}
			return v, nil
		}),
		rxgo.CatchError2[int](func(err error) rxgo.Observable[string] {
			return rxgo.Of([]string{"I", "II", "III", "IV", "V"})
		}),
	), []rxgo.Either[int, string]{
		rxgo.NewA[int, string](2),
		rxgo.NewA[int, string](4),
		rxgo.NewA[int, string](6),
		rxgo.NewA[int, string](8),
		rxgo.NewA[int, string](10),
		rxgo.NewB[int]("I"),
		rxgo.NewB[int]("II"),
		rxgo.NewB[int]("III"),
		rxgo.NewB[int]("IV"),
		rxgo.NewB[int]("V"),
	})
}

func TestDefaultIfEmpty(t *testing.T) {
	defer goleak.VerifyNone(t)

	assertItem(t, rxgo.Pipe1(
		rxgo.Empty[string](),
		rxgo.DefaultIfEmpty(`hello world!`),
	), []string{`hello world!`})
}

func TestDistinctUntilChanged(t *testing.T) {
	defer goleak.VerifyNone(t)

	assertItem(t, rxgo.Pipe1(
		rxgo.Of([]int{1, 1, 1, 2, 2, 2, 1, 1, 3, 3}),
		rxgo.DistinctUntilChanged[int](),
	), []int{1, 2, 1, 3})
}

func TestEvery(t *testing.T) {
	defer goleak.VerifyNone(t)

	t.Run("true", func(t *testing.T) {
		assertItem(t, rxgo.Pipe1(
			rxgo.Of([]int{1, 1, 1, 1, 1, 3, 3}),
			rxgo.Every(func(v int, _ int) bool {
				return v%2 != 0
			}),
		), []bool{true})
	})

	t.Run("false", func(t *testing.T) {
		assertItem(t, rxgo.Pipe1(
			rxgo.Of([]int{1, 1, 1, 1, 1, 3, 3}),
			rxgo.Every(func(v int, _ int) bool {
				return v%2 == 0
			}),
		), []bool{false})
	})
}

func TestFilter(t *testing.T) {
	defer goleak.VerifyNone(t)

	assertItem(t, rxgo.Pipe1(
		rxgo.Of([]int{1, 2, 3, 4, 7, 8, 9, 10, 11, 14}),
		rxgo.Filter(func(v int) bool {
			return v%2 == 0
		}),
	), []int{2, 4, 8, 10, 14})
}

func TestIsEmpty(t *testing.T) {
	defer goleak.VerifyNone(t)

	t.Run("Empty", func(t *testing.T) {
		assertItem(t, rxgo.Pipe1(
			rxgo.Empty[any](),
			rxgo.IsEmpty[any](),
		), []bool{true})
	})

	t.Run("Not empty", func(t *testing.T) {
		assertItem(t, rxgo.Pipe1(
			rxgo.Of([]int{1}),
			rxgo.IsEmpty[int](),
		), []bool{false})
	})
}

func TestRange(t *testing.T) {
	defer goleak.VerifyNone(t)

	t.Run("Negative range values", func(t *testing.T) {
		assertItem(t, rxgo.Range(-5, 1), []int{-5, -4, -3, -2, -1, 0, 1})
	})
	t.Run("Positive range values", func(t *testing.T) {
		assertItem(t, rxgo.Range(1, 5), []int{1, 2, 3, 4, 5})
	})
}

func TestReduce(t *testing.T) {
	defer goleak.VerifyNone(t)

	t.Run("Negative values", func(t *testing.T) {
		assertItem(t, rxgo.Pipe1(
			rxgo.Of([]int{-88, 1, 2, -3, 888}),
			rxgo.Reduce(func(acc int, v int, _ int) int {
				return acc + v
			}, int(0)),
		), []int{800})
	})

	t.Run("Positive values", func(t *testing.T) {
		assertItem(t, rxgo.Pipe1(
			rxgo.Of([]uint{100, 1, 2, 888}),
			rxgo.Reduce(func(acc uint, v uint, _ int) uint {
				return acc + v
			}, uint(0)),
		), []uint{991})
	})
}

func TestSkipLast(t *testing.T) {
	defer goleak.VerifyNone(t)

	assertItem(t, rxgo.Pipe1(
		rxgo.Of([]string{"a", "b", "c", "d", "e", "f"}),
		rxgo.SkipLast[string](3),
	), []string{"a", "b", "c"})
}

func TestScan(t *testing.T) {
	defer goleak.VerifyNone(t)

	assertItem(t, rxgo.Pipe1(
		rxgo.Of([]int{1, 2, 3}),
		rxgo.Scan(func(acc int, v int, _ int) int {
			return acc + v
		}, 0)), []int{1, 3, 6})
}

func TestTakeLast(t *testing.T) {
	defer goleak.VerifyNone(t)

	assertItem(t, rxgo.Pipe1(
		rxgo.Of([]string{"a", "b", "c", "d", "e", "f"}),
		rxgo.TakeLast[string](1),
	), []string{"f"})
}

func TestMin(t *testing.T) {
	defer goleak.VerifyNone(t)

	assertItem(t, rxgo.Pipe1(
		rxgo.Of([]int{-88, 1, 2, 3}),
		rxgo.Min[int](),
	), []int{-88})
}

func TestMax(t *testing.T) {
	defer goleak.VerifyNone(t)

	assertItem(t, rxgo.Pipe1(
		rxgo.Of([]int{-88, 1, 2, 3, 888}),
		rxgo.Max[int](),
	), []int{888})
}

func assertItem[T any](t *testing.T, observable rxgo.Observable[T], expected []T) {
	result := make([]T, 0)
	for v, err := range observable.Subscribe() {
		require.NoError(t, err)
		result = append(result, v)
	}
	require.ElementsMatch(t, result, expected)
}

func assertItems[T any](t *testing.T, observable rxgo.Observable[[]T], expected [][]T) {
	result := make([][]T, 0)
	for v, err := range observable.Subscribe() {
		require.NoError(t, err)
		result = append(result, append([]T{}, v...))
	}
	require.ElementsMatch(t, result, expected)
}
