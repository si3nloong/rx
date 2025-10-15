# IgnoreElements

> Ignores all items emitted by the source Observable and only passes calls of complete or error.

## Description

![](https://rxjs.dev/assets/images/marble-diagrams/ignoreElements.png)

The `IgnoreElements` operator suppresses all items emitted by the source Observable, but allows its termination notification (either error or complete) to pass through unchanged.

If you do not care about the items being emitted by an Observable, but you do want to be notified when it completes or when it terminates with an error, you can apply the `IgnoreElements` operator to the Observable, which will ensure that it will never call its observers’ next handlers.

## Example

```go
for v, _ := range rxgo.Pipe1(
    rxgo.Of("You", "talking", "to", "me"),
    rxgo.IgnoreElements(),
).Subscribe() {
    println(v)
}
println("the end")
```

Output:

```
the end
```
