# Empty

> A simple Observable that emits no items to the Observer and immediately emits a complete notification.

## Description

> Just emits 'complete', and nothing else.

![](https://rxjs.dev/assets/images/marble-diagrams/empty.png)

A simple Observable that only emits the complete notification. It can be used for composing with other Observables, such as in a `MergeMap`.

## Example

```go
for v, _ := range rx.Empty[string]().Subscribe() {
    println(v)
}
println("Completed!")
```

Output:

```
Completed!
```
