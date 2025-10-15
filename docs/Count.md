# Count

> Counts the number of emissions on the source and emits that number when the source completes.

## Description

> Tells how many values were emitted, when the source completes.

![](https://rxjs.dev/assets/images/marble-diagrams/count.png)

`Count` transforms an Observable that emits values into an Observable that emits a single value that represents the number of values emitted by the source Observable. If the source Observable terminates with an error, `Count` will pass this error notification along without emitting a value first. If the source Observable does not terminate at all, count will neither emit a value nor terminate. This operator takes an optional `predicate` function as argument, in which case the output emission will represent the number of source values that matched `true` with the `predicate`.

## Example

```go
for v, err := range rxgo.Pipe1(
    rxgo.Of(1, 2, 3, 4),
    rxgo.Count(),
).Subscribe() {
    if err != nil {
        panic(err)
    }
    println(v)
}
```

Output:

```
4
```
