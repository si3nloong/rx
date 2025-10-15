# ElementAt

> Emits the single value at the specified index in a sequence of emissions from the source Observable.

## Description

> Emits only the i-th value, then completes.

![](https://rxjs.dev/assets/images/marble-diagrams/elementAt.png)

`ElementAt` returns an Observable that emits the item at the specified index in the source Observable, or a default value if that index is out of range and the default argument is provided. If the default argument is not given and the index is out of range, the output Observable will emit an `ErrArgumentOutOfRange` error.

## Example

```go
for v, err := range rxgo.Pipe1(
    rxgo.Of(1, 2, 3, 4, 5),
    rxgo.ElementAt[int](2),
).Subscribe() {
    if err != nil {
        panic(err)
    }
    println(v)
}
```

Output:

```
3
```
