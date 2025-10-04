# Defer

> Creates an Observable that, on subscribe, calls an Observable factory to make an Observable for each new Observer.

## Description

> Creates the Observable lazily, that is, only when it is subscribed.

![](https://rxjs.dev/assets/images/marble-diagrams/defer.png)

`Defer` allows you to create an Observable only when the Observer subscribes. It waits until an Observer subscribes to it, calls the given factory function to get an Observable -- where a factory function typically generates a new Observable -- and subscribes the Observer to this Observable. In case the factory function returns a falsy value, then `Empty` is used as Observable instead. Last but not least, an exception during the factory function call is transferred to the Observer by calling error.

## Example

```go
for v, _ := range rxgo.Defer(func() rxgo.Observable[int] {
    return rxgo.ObservableFunc[int](func(yield func(int, error) bool) {
        for i := 0; i < 10; i++ {
            if !yield(i, nil) {
                return
            }
        }
    })
}) {
    println(v)
}
```

Output:

```
0
1
2
3
4
5
6
7
8
9
```
