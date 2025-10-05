package rxgo

type Kind uint8

const (
	KindNext Kind = iota
	KindError
	KindComplete
)

type Notification[T any] struct {
	kind Kind
	err  error
	v    T
}

func (n Notification[T]) Kind() Kind {
	return n.kind
}
func (n *Notification[T]) Value() T {
	return n.v
}
func (n *Notification[T]) Error() error {
	return n.err
}

func NewNextNotification[T any](v T) *Notification[T] {
	return &Notification[T]{
		kind: KindNext,
		v:    v,
	}
}

func NewErrorNotification[T any](v T) *Notification[T] {
	return &Notification[T]{
		kind: KindNext,
		v:    v,
	}
}

func NewCompleteNotification[T any]() *Notification[T] {
	return &Notification[T]{
		kind: KindComplete,
	}
}
