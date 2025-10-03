# RxGo

Reactive Extensions for the Go Language.

## ReactiveX

[ReactiveX](http://reactivex.io/), or Rx for short, is an API for programming with Observable streams. This is the official ReactiveX API for the Go language.

ReactiveX is a new, alternative way of asynchronous programming to callbacks, promises, and deferred. It is about processing streams of events or items, with events being any occurrences or changes within the system. A stream of events is called an [Observable](http://reactivex.io/documentation/contract.html).

An operator is a function that defines an Observable, how and when it should emit data. The list of operators covered is available [here](README.md#).

## Why another RxGo?

- The official [RxGo](https://github.com/ReactiveX/RxGo) library is not maintainable anymore.
- Go generics by default.
- Utilise iterator pattern from Go which available since [Go 1.23](https://go.dev/blog/go1.23).
- **Zero dependencies**, standard library only.


## Installation

```go
go get -u github.com/si3nloong/rxgo
```

## Getting Started

There is no magic under the hood, an observable is just a [Go iterator](https://go.dev/blog/range-functions) which comply to `iter.Seq2[T, error]`.

You can create an Observable as easy as :

```go
rxgo.ObservableFunc[string](func(yield func(string, error) bool) {
	if !yield("hello", nil) {
		return
	}
})
```

```go
package main

import (
	"log"

	"github.com/si3nloong/rxgo"
)

func main() {
    rxgo.Pipe2(
		rxgo.Of([]int{1, 1, 1, 2, 2, 2, 1, 1, 3, 3}),
		rxgo.DistinctUntilChanged[int](),
		rxgo.ToSlice[int](),
	).Subscribe(func(v []int) {
		log.Println(v)
	}, func(err error) {}, func() {
		println("Completed!")
	})
    // 1, 2, 1, 3, Completed!
}
```

## Advanced

```go
    rxgo.Pipe3(
        // Emit every one second
		rxgo.Interval(time.Second),
		rxgo.Filter(func(v int) bool {
            // Filter with modules
			return v%2 == 0
		}),
		rxgo.Map2(func(v int, _ int) (int, error) {
			if v > 10 {
                // Throw error when value is greather than 10
				return 0, errors.New(`stop la`)
			}
			return v, nil
		}),
        // Catch error and return new Observable
		rxgo.CatchError2[int](func(err error) rxgo.Observable[string] {
			return rxgo.Of([]string{"I", "II", "III", "IV", "V"})
		}),
	).SubscribeOn(func(v rxgo.Either[int, string]) {
		if v1, ok := v.A(); ok {
			log.Println("Result ->", v1)
		} else {
			log.Println("Result ->", v.MustB())
		}
	}, func(err error) {
		log.Println("Error ->", err)
	}, func() {
        log.Println("Completed!")
    })
    // 0, 2, 4, 6, 8, 10, I, II, III, IV, V, Completed!
```