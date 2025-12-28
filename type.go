package rx

import "time"

// Either represents a value of one of two possible types (a disjoint union).
type Either[A, B any] struct {
	v any
}

// NewA creates an Either with the left value.
func NewA[B, A any](v A) Either[A, B] {
	return Either[A, B]{v: v}
}

// NewB creates an Either with the right value.
func NewB[A, B any](v B) Either[A, B] {
	return Either[A, B]{v: v}
}

// A returns the left value and a boolean indicating if it exists.
func (a Either[A, B]) A() (A, bool) {
	if v, ok := a.v.(A); ok {
		return v, true
	}
	var zero A
	return zero, false
}

// B returns the right value and a boolean indicating if it exists.
func (a Either[A, B]) B() (B, bool) {
	if v, ok := a.v.(B); ok {
		return v, true
	}
	var zero B
	return zero, false
}

// MustA returns the left value. It panics if the value is not of type A.
func (a Either[A, B]) MustA() A {
	return a.v.(A)
}

// MustB returns the right value. It panics if the value is not of type B.
func (a Either[A, B]) MustB() B {
	return a.v.(B)
}

// Tuple represents a pair of values.
type Tuple[L, R any] struct {
	Left  L
	Right R
}

// NewTuple creates a new Tuple.
func NewTuple[L, R any](left L, right R) Tuple[L, R] {
	return Tuple[L, R]{left, right}
}

// TimeInterval represents a value emitted by an Observable with the time interval since the last emission.
type TimeInterval[T any] struct {
	Interval time.Duration
	Value    T
}

// Timestamp represents a value emitted by an Observable with the timestamp of the emission.
type Timestamp[T any] struct {
	Time  time.Time
	Value T
}
