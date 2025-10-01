# RxGo

Reactive Extensions for the Go Language.

## ReactiveX

[ReactiveX](http://reactivex.io/), or Rx for short, is an API for programming with Observable streams. This is the official ReactiveX API for the Go language.

ReactiveX is a new, alternative way of asynchronous programming to callbacks, promises, and deferred. It is about processing streams of events or items, with events being any occurrences or changes within the system. A stream of events is called an [Observable](http://reactivex.io/documentation/contract.html).

An operator is a function that defines an Observable, how and when it should emit data. The list of operators covered is available [here](README.md#supported-operators-in-rxgo).

## Why another RxGo?

- The official [RxGo](https://github.com/ReactiveX/RxGo) library. is not maintainable anymore.
- Support Go generics by default.
- Use iterator pattern from Go which available since [Go 1.23](https://go.dev/blog/go1.23).


## Examples

```go
package main

import (
	"log"
	"math/rand/v2"
	"time"

	"github.com/si3nloong/rxgo"
)

func main() {
    rxgo.Pipe4(
		rxgo.Interval(time.Second),
		rxgo.Tap(func(v int) {
			println(v)
		}),
		rxgo.Buffer[int](func(yield func(int, error) bool) {
			const (
				min = 2
				max = 10
			)

			for {
				d := time.Second * time.Duration(rand.IntN(max-min)+min)
				println(d.String())
				<-time.After(d)
				if !yield(0, nil) {
					return
				}
			}
		}),
		rxgo.Take[[]int](3),
		rxgo.Delay[[]int](time.Second),
	).Subscribe(func(v []int) {
		log.Println("Next ->", v)
	}, func(err error) {
		log.Println("Error ->", err)
	}, func() {
        log.Println("Completed!")
    })
}
```