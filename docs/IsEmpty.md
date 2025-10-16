# IsEmpty

> Emits `false` if the input Observable emits any values, or emits true if the input Observable completes without emitting any values.

## Description

> Tells whether any values are emitted by an Observable.

![](https://rxjs.dev/assets/images/marble-diagrams/isEmpty.png)

`IsEmpty` transforms an Observable that emits values into an Observable that emits a single boolean value representing whether or not any values were emitted by the source Observable. As soon as the source Observable emits a value, `IsEmpty` will emit a `false` and complete. If the source Observable completes having not emitted anything, `IsEmpty` will emit a `true` and complete.

A similar effect could be achieved with `Count`, but isEmpty can emit a `false` value sooner.

## Example

```go
for v, _ := range rx.Pipe1(
    rx.Empty[string](),
    rx.IsEmpty(),
).Subscribe() {
    println(v)
}
```

Output:

```
true
```
