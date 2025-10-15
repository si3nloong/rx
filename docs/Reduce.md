# Reduce

> Applies an accumulator function over the source Observable, and returns the accumulated result when the source completes, given an optional seed value.

## Description

> Combines together all values emitted on the source, using an accumulator function that knows how to join a new source value into the accumulation from the past.

![](https://rxjs.dev/assets/images/marble-diagrams/reduce.png)

`Reduce` applies an accumulator function against an accumulation and each value of the source Observable (from the past) to reduce it to a single value, emitted on the output Observable. Note that `Reduce` will only emit one value, only when the source Observable completes. It is equivalent to applying operator `Scan` followed by operator `Last`.

Returns an Observable that applies a specified accumulator function to each item emitted by the source Observable. If a seed value is specified, then that value will be used as the initial value for the accumulator. If no seed value is specified, the first item of the source is used as the seed.

## Example

```go
for v, _ := range rxgo.Pipe1(
    rxgo.Of(1, 3, 4),
    rxgo.Reduce(func(acc int, v int, _ int) int {
        return acc + v
    }, 0),
).Subscribe() {
    println(v)
}
```

Output:

```
8
```
