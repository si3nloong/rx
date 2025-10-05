package rxgo

import (
	"context"
	"iter"
	"sync"
)

type Subject[T any] struct {
	ch             chan Notification[T]
	broadcast      chan Notification[T]
	closed         bool
	emitOnce       sync.Once
	closeOnce      sync.Once
	listeners      []chan Notification[T]
	addListener    chan chan Notification[T]
	removeListener chan chan Notification[T]
	stopCh         chan struct{}
}

func NewSubject[T any]() *Subject[T] {
	subject := &Subject[T]{
		ch:             make(chan Notification[T], 1),
		broadcast:      make(chan Notification[T], 1),
		stopCh:         make(chan struct{}, 1),
		listeners:      make([]chan Notification[T], 0),
		addListener:    make(chan chan Notification[T], 1),
		removeListener: make(chan chan Notification[T], 1),
	}
	go subject.run(context.Background())
	return subject
}

func (s *Subject[T]) run(ctx context.Context) {
	defer s.Close()

	for {
		select {
		case <-ctx.Done():
			return
		case v := <-s.addListener:
			s.listeners = append(s.listeners, v)
		case v := <-s.removeListener:
			for i, ch := range s.listeners {
				if ch == v {
					s.listeners[i] = s.listeners[len(s.listeners)-1]
					s.listeners = s.listeners[:len(s.listeners)-1]
					close(ch)
					break
				}
			}
		case v, ok := <-s.ch:
			if !ok {
				return
			}
			for _, listener := range s.listeners {
				if listener != nil {
					select {
					case <-ctx.Done():
						return
					case listener <- v:
					}
				}
			}
		}
	}
}

func (s *Subject[T]) Next(v T) {
	if s.closed {
		return
	}
	s.ch <- Notification[T]{kind: KindNext, v: v}
}

func (s *Subject[T]) Error(err error) {
	s.emitOnce.Do(func() {
		s.ch <- Notification[T]{kind: KindError, err: err}
	})
}

func (s *Subject[T]) Complete() {
	s.emitOnce.Do(func() {
		s.ch <- Notification[T]{kind: KindComplete}
	})
}

func (s *Subject[T]) Close() {
	s.closeOnce.Do(func() {
		s.closed = true
		for _, listener := range s.listeners {
			if listener != nil {
				// close(listener)
			}
		}
		close(s.ch)
		close(s.addListener)
		close(s.removeListener)
		close(s.broadcast)
		close(s.stopCh)
	})
}

func (s *Subject[T]) Subscribe() iter.Seq2[T, error] {
	return func(yield func(T, error) bool) {
		subscription := make(chan Notification[T])
		s.addListener <- subscription
		defer close(subscription)

		for {
			select {
			case v, ok := <-subscription:
				if !ok {
					return
				}
				switch v.Kind() {
				case KindError:
					var zero T
					yield(zero, v.Error())
					return
				case KindComplete:
					return
				case KindNext:
					if !yield(v.Value(), nil) {
						return
					}
				default:
					panic("unreachable")
				}
			default:
			}
		}
	}
}
