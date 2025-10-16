# First

> Emits only the first value (or the first value that meets some condition) emitted by the source Observable.

## Description

> Emits only the first value. Or emits only the first value that passes some test.

![](https://rxjs.dev/assets/images/marble-diagrams/first.png)

If called with no arguments, first emits the first value of the source Observable, then completes. If called with a predicate function, first emits the first value of the source that matches the specified condition. Emits an error notification if defaultValue was not provided and a matching element is not found.

## Example

```go
for v, _ := range rx.Pipe1(
    rx.From[string]([]string{"Ben", "Tracy", "Laney", "Lily"}),
    rx.First[string](),
).Subscribe() {
    println(v)
}
```

Output:

```
"Ben"
```
