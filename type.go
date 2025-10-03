package rxgo

import "time"

type Either[A, B any] struct {
	v any
}

func (a Either[A, B]) A() (A, bool) {
	if v, ok := a.v.(A); ok {
		return v, true
	}
	var zero A
	return zero, false
}

func (a Either[A, B]) B() (B, bool) {
	if v, ok := a.v.(B); ok {
		return v, true
	}
	var zero B
	return zero, false
}

func (a Either[A, B]) MustA() A {
	return a.v.(A)
}
func (a Either[A, B]) MustB() B {
	return a.v.(B)
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
