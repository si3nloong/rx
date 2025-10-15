# DefaultIfEmpty

> Emits a given value if the source Observable completes without emitting any next value, otherwise mirrors the source Observable.

## Description

> If the source Observable turns out to be empty, then this operator will emit a default value.

![](https://rxjs.dev/assets/images/marble-diagrams/defaultIfEmpty.png)


## Example

```go
for v, err := range rxgo.Pipe1(
    rxgo.Empty[string](),
    rxgo.DefaultIfEmpty(`hello world!`),
).Subscribe() {
    if err != nil {
        panic(err)
    }
    println(v)
}
```

Output:

```
hello world!
```
