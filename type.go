package rx

import "time"

type Either[A, B any] struct {
	v any
}

func NewA[B, A any](v A) Either[A, B] {
	return Either[A, B]{v: v}
}

func NewB[A, B any](v B) Either[A, B] {
	return Either[A, B]{v: v}
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

type Tuple[L, R any] struct {
	Left  L
	Right R
}

func NewTuple[L, R any](left L, right R) Tuple[L, R] {
	return Tuple[L, R]{left, right}
}

type TimeInterval[T any] struct {
	Interval time.Duration
	Value    T
}

type Timestamp[T any] struct {
	Time  time.Time
	Value T
}
