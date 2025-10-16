package main_test

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/si3nloong/rx"
	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

func TestBuffer(t *testing.T) {
	defer goleak.VerifyNone(t)

	assertItems(t, rx.Pipe1(
		rx.From[int]([]int{1, 3, 4, 5, 9}),
		rx.BufferCount[int](2),
	), [][]int{
		{1, 3},
		{4, 5},
		{9},
	})
}

func TestCombineLatest(t *testing.T) {
	defer goleak.VerifyNone(t)

	assertItems(t, rx.CombineLatest(
		rx.Pipe2(rx.From[int]([]int{1}), rx.Delay[int](time.Second), rx.StartWith(0)),
		rx.Pipe2(rx.From[int]([]int{5}), rx.Delay[int](time.Second*5), rx.StartWith(0)),
		rx.Pipe2(rx.From[int]([]int{10}), rx.Delay[int](time.Second*10), rx.StartWith(0)),
	), [][]int{
		{0, 0, 0},
		{1, 0, 0},
		{1, 5, 0},
		{1, 5, 10},
	})
}

func TestCatchError(t *testing.T) {
	defer goleak.VerifyNone(t)

	assertItem(t, rx.Pipe3(
		rx.Range(1, 20),
		rx.Filter(func(v int) bool {
			return v%2 == 0
		}),
		rx.Map2(func(v int, _ int) (int, error) {
			if v > 10 {
				return 0, errors.New(`stop la`)
			}
			return v, nil
		}),
		rx.CatchError2[int](func(err error) rx.Observable[string] {
			return rx.From[string]([]string{"I", "II", "III", "IV", "V"})
		}),
	), []rx.Either[int, string]{
		rx.NewA[string](2),
		rx.NewA[string](4),
		rx.NewA[string](6),
		rx.NewA[string](8),
		rx.NewA[string](10),
		rx.NewB[int]("I"),
		rx.NewB[int]("II"),
		rx.NewB[int]("III"),
		rx.NewB[int]("IV"),
		rx.NewB[int]("V"),
	})
}

func TestDefaultIfEmpty(t *testing.T) {
	defer goleak.VerifyNone(t)

	assertItem(t, rx.Pipe1(
		rx.Empty[string](),
		rx.DefaultIfEmpty(`hello world!`),
	), []string{`hello world!`})
}

func TestDistinct(t *testing.T) {
	defer goleak.VerifyNone(t)

	assertItem(t, rx.Pipe1(
		rx.From[int]([]int{1, 1, 2, 2, 2, 1, 2, 3, 4, 3, 2, 1}),
		rx.Distinct(func(v int) int {
			return v
		}),
	), []int{1, 2, 3, 4})

	type user struct {
		name string
		age  int
	}

	assertItem(t, rx.Pipe1(
		rx.From[user]([]user{
			{"Foo", 4},
			{"Bar", 7},
			{"Foo", 5},
		}),
		rx.Distinct(func(v user) string {
			return v.name
		}),
	), []user{{"Foo", 4}, {"Bar", 7}})
}

func TestDistinctUntilChanged(t *testing.T) {
	defer goleak.VerifyNone(t)

	assertItem(t, rx.Pipe1(
		rx.Of(1, 1, 1, 2, 2, 2, 1, 1, 3, 3),
		rx.DistinctUntilChanged[int](),
	), []int{1, 2, 1, 3})
}

func TestElementAt(t *testing.T) {
	defer goleak.VerifyNone(t)

	t.Run("Valid index", func(t *testing.T) {
		assertItem(t, rx.Pipe1(
			rx.Of(1, 2, 3, 4, 5),
			rx.ElementAt[int](2),
		), []int{3})
	})

	t.Run("Invalid index", func(t *testing.T) {
		isError(t, rx.Pipe1(
			rx.Of(1, 2, 3, 4, 5),
			rx.ElementAt[int](10),
		), rx.ErrArgumentOutOfRange)
	})
}

func TestEvery(t *testing.T) {
	defer goleak.VerifyNone(t)

	t.Run("true", func(t *testing.T) {
		assertItem(t, rx.Pipe1(
			rx.Of(1, 1, 1, 1, 1, 3, 3),
			rx.Every(func(v int, _ int) bool {
				return v%2 != 0
			}),
		), []bool{true})
	})

	t.Run("false", func(t *testing.T) {
		assertItem(t, rx.Pipe1(
			rx.Of(1, 1, 1, 1, 1, 3, 3),
			rx.Every(func(v int, _ int) bool {
				return v%2 == 0
			}),
		), []bool{false})
	})
}

func TestFirst(t *testing.T) {
	defer goleak.VerifyNone(t)

	assertItem(t, rx.Pipe1(
		rx.From[string]([]string{"Ben", "Tracy", "Laney", "Lily"}),
		rx.First[string](),
	), []string{"Ben"})
}

func TestFilter(t *testing.T) {
	defer goleak.VerifyNone(t)

	assertItem(t, rx.Pipe1(
		rx.Of(1, 2, 3, 4, 7, 8, 9, 10, 11, 14),
		rx.Filter(func(v int) bool {
			return v%2 == 0
		}),
	), []int{2, 4, 8, 10, 14})
}

func TestForkJoin(t *testing.T) {
	defer goleak.VerifyNone(t)

	assertItems(t, rx.ForkJoin(
		rx.Of(1, 2, 3, 4, 7, 8, 9, 10, 11, 14),
		rx.Timer[int](time.Second),
	), [][]int{{14, 0}})

	assertItem(t, rx.ForkJoin(
		rx.Empty[int](),
		rx.Of(1, 2, 3, 4, 7, 8, 9, 10, 11, 14),
		rx.Timer[int](time.Second),
	), [][]int{})
}

func TestIsEmpty(t *testing.T) {
	defer goleak.VerifyNone(t)

	t.Run("Empty", func(t *testing.T) {
		assertItem(t, rx.Pipe1(
			rx.Empty[any](),
			rx.IsEmpty[any](),
		), []bool{true})
	})

	t.Run("Not empty", func(t *testing.T) {
		assertItem(t, rx.Pipe1(
			rx.Of(1),
			rx.IsEmpty[int](),
		), []bool{false})
	})
}

func TestRange(t *testing.T) {
	defer goleak.VerifyNone(t)

	t.Run("Negative range values", func(t *testing.T) {
		assertItem(t, rx.Range(-5, 1), []int{-5, -4, -3, -2, -1, 0, 1})
	})
	t.Run("Positive range values", func(t *testing.T) {
		assertItem(t, rx.Range(1, 5), []int{1, 2, 3, 4, 5})
	})
}

func TestReduce(t *testing.T) {
	defer goleak.VerifyNone(t)

	t.Run("Negative values", func(t *testing.T) {
		assertItem(t, rx.Pipe1(
			rx.Of(-88, 1, 2, -3, 888),
			rx.Reduce(func(acc int, v int, _ int) int {
				return acc + v
			}, int(0)),
		), []int{800})
	})

	t.Run("Positive values", func(t *testing.T) {
		assertItem(t, rx.Pipe1(
			rx.From[uint]([]uint{100, 1, 2, 888}),
			rx.Reduce(func(acc uint, v uint, _ int) uint {
				return acc + v
			}, uint(0)),
		), []uint{991})
	})
}

func TestSingle(t *testing.T) {
	defer goleak.VerifyNone(t)

	assertItem(t, rx.Pipe1(
		rx.From[string]([]string{"Ben", "Tracy", "Laney", "Lily"}),
		rx.Single(func(v string, _ int) bool {
			return strings.HasPrefix(v, "B")
		}),
	), []string{"Ben"})

	t.Run("Empty", func(t *testing.T) {
		isError(t, rx.Pipe1(
			rx.From[string]([]string{}),
			rx.Single(func(v string, _ int) bool {
				return strings.HasPrefix(v, "B")
			}),
		), rx.ErrEmpty)
	})

	t.Run("Too many value matched", func(t *testing.T) {
		isError(t, rx.Pipe1(
			rx.From[string]([]string{"Ben", "Tracy", "Bradley", "Lily"}),
			rx.Single(func(v string, _ int) bool {
				return strings.HasPrefix(v, "B")
			}),
		), rx.ErrSequence)
	})

	t.Run("Not found", func(t *testing.T) {
		isError(t, rx.Pipe1(
			rx.From[string]([]string{"Lance", "Lily"}),
			rx.Single(func(v string, _ int) bool {
				return strings.HasPrefix(v, "B")
			}),
		), rx.ErrNotFound)
	})
}

func TestSkipLast(t *testing.T) {
	defer goleak.VerifyNone(t)

	assertItem(t, rx.Pipe1(
		rx.Of("a", "b", "c", "d", "e", "f"),
		rx.SkipLast[string](3),
	), []string{"a", "b", "c"})
}

func TestScan(t *testing.T) {
	defer goleak.VerifyNone(t)

	assertItem(t, rx.Pipe1(
		rx.Of(1, 2, 3),
		rx.Scan(func(acc int, v int, _ int) int {
			return acc + v
		}, 0)), []int{1, 3, 6})
}

func TestTakeLast(t *testing.T) {
	defer goleak.VerifyNone(t)

	assertItem(t, rx.Pipe1(
		rx.From[string]([]string{"a", "b", "c", "d", "e", "f"}),
		rx.TakeLast[string](1),
	), []string{"f"})
}

func TestMin(t *testing.T) {
	defer goleak.VerifyNone(t)

	assertItem(t, rx.Pipe1(
		rx.Of(-88, 1, 2, 3),
		rx.Min[int](),
	), []int{-88})
}

func TestMax(t *testing.T) {
	defer goleak.VerifyNone(t)

	assertItem(t, rx.Pipe1(
		rx.Of(-88, 1, 2, 3, 888),
		rx.Max[int](),
	), []int{888})
}

func assertItem[T any](t *testing.T, observable rx.Observable[T], expected []T) {
	result := make([]T, 0)
	for v, err := range observable.Subscribe() {
		require.NoError(t, err)
		result = append(result, v)
	}
	require.ElementsMatch(t, result, expected)
}

func assertItems[T any](t *testing.T, observable rx.Observable[[]T], expected [][]T) {
	result := make([][]T, 0)
	for v, err := range observable.Subscribe() {
		require.NoError(t, err)
		result = append(result, append([]T{}, v...))
	}
	require.ElementsMatch(t, result, expected)
}

func isError[T any](t *testing.T, observable rx.Observable[T], expectedErr error) {
	for _, err := range observable.Subscribe() {
		if err != nil {
			require.ErrorIs(t, err, expectedErr)
			return
		}
		require.NoError(t, err)
	}
}
