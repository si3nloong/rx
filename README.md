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
- No extra dependencies, everything using standard lib.


## Examples

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