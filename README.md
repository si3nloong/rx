# RxGo

Reactive Extensions for the Go Language.

## ReactiveX

[ReactiveX](http://reactivex.io/), or Rx for short, is an API for programming with Observable streams. This is the official ReactiveX API for the Go language.

ReactiveX is a new, alternative way of asynchronous programming to callbacks, promises, and deferred. It is about processing streams of events or items, with events being any occurrences or changes within the system. A stream of events is called an [Observable](http://reactivex.io/documentation/contract.html).

An operator is a function that defines an Observable, how and when it should emit data. The list of operators covered is available [here](README.md#).

## Why another RxGo?

- The official [RxGo](https://github.com/ReactiveX/RxGo) library is no longer maintainable anymore.
- Go generics by default, no reflection.
- Utilise iterator pattern from Go which available since [Go 1.23](https://go.dev/blog/go1.23).
- **Zero dependencies**, standard library only.

## Installation

```go
go get -u github.com/si3nloong/rx
```

## Getting Started

There is no magic under the hood, an observable is just a [Go iterator](https://go.dev/blog/range-functions) which comply to [iter.Seq2[T, error]](https://pkg.go.dev/iter#Seq2) interface.

You can create an Observable as easy as :

```go
observable := rx.ObservableFunc[string](func(yield func(string, error) bool) {
	if !yield("hello", nil) {
		return
	}
})
```

Once the Observable is created, we can observe it using the `Subscribe` or `SubscribeOn` function. By default, an Observable is lazy in the sense that it emits items only once a subscription is made.

Since it is a Go iterator, you can observe it using push method with `range` keyword:

```go
for v, err := range observable.Subscribe() {
	if err != nil {
		panic(err)
	}
	println(v)
}
```

OR we can observe it using pull method with [iter.Pull2](https://pkg.go.dev/iter#Pull2) function:

```go
next, stop := iter.Pull2(observable.Subscribe())
defer stop()

for {
	v, err, ok := next()
	if err != nil {
		panic(err)
	} else if !ok {
		println("Completed!")
		return
	} else {
		println(v)
	}
}
```

## Categories of operators

There are operators for different purposes, and they may be categorized as: creation, transformation, filtering, joining, multicasting, error handling, utility, etc.

For a complete overview, see the [operators](/docs/operators/README.md) page.

## Observable Types

### Hot vs. Cold Observables

In the Rx world, there is a distinction between cold and hot Observables. When the data is produced by the Observable itself, it is a cold Observable. When the data is produced outside the Observable, it is a hot Observable. Usually, when we don't want to create a producer over and over again, we favour a hot Observable.

In RxGo, there is a similar concept.

## Performance

![](/benchmarks/result.png)

## Example

```go
package main

import (
	"github.com/si3nloong/rx"
)

func main() {
    for v, err := rx.Pipe1(
		rx.Of(1, 1, 1, 2, 2, 2, 1, 1, 3, 3),
		rx.DistinctUntilChanged[int](),
	) {
		if err != nil {
			panic(err)
		}
		println(v)
	}
	println("Completed!")
    // 1, 2, 1, 3, Completed!
}
```
