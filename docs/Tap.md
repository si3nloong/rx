# Tap

> Used to perform side-effects for notifications from the source observable

## Description

> Used when you want to affect outside state with a notification without altering the notification

![](https://rxjs.dev/assets/images/marble-diagrams/tap.png)

`Tap` is designed to allow the developer a designated place to perform side effects. While you could perform side-effects inside of a `Map` or a `MergeMap`, that would make their mapping functions impure, which isn't always a big deal, but will make it so you can't do things like memoize those functions. The `Tap` operator is designed solely for such side-effects to help you remove side-effects from other operations.

For any notification, next, error, or complete, `Tap` will call the appropriate callback you have provided to it, via a function reference, or a partial observer, then pass that notification down the stream.

The observable returned by `Tap` is an exact mirror of the source, with one exception: Any error that occurs -- synchronously -- in a handler provided to `Tap` will be emitted as an error from the returned observable.

Be careful! You can mutate objects as they pass through the `Tap` operator's handlers.

The most common use of `Tap` is actually for debugging.

## Example

```go
for v, _ := range rxgo.Pipe1(
    rxgo.Of(1, 3, 4),
    rxgo.ToSlice(),
) {
    fmt.Println(v)
}
```

Output:

```
[1, 3, 4]
```