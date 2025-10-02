package rxgo

import "time"

type Either[T1, T2 any] struct {
	v any
}

func (a Either[T1, T2]) Type1() (T1, bool) {
	if v, ok := a.v.(T1); ok {
		return v, true
	}
	var zero T1
	return zero, false
}

func (a Either[T1, T2]) Type2() (T2, bool) {
	if v, ok := a.v.(T2); ok {
		return v, true
	}
	var zero T2
	return zero, false
}

func (a Either[T1, T2]) MustType1() T1 {
	return a.v.(T1)
}
func (a Either[T1, T2]) MustType2() T2 {
	return a.v.(T2)
}

type Tuple[A, B any] struct {
	a A
	b B
}

func (t Tuple[A, B]) Left() A {
	return t.a
}
func (t Tuple[A, B]) Right() B {
	return t.b
}

type TimeInterval[T any] struct {
	Interval time.Duration
	Value    T
}

type Timestamp[T any] struct {
	Time  time.Time
	Value T
}
