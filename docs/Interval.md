# Interval

> Creates an Observable that emits sequential numbers every specified interval of time.

## Description

> Emits incremental numbers periodically in time.

![](https://rxjs.dev/assets/images/marble-diagrams/interval.png)

`Interval` returns an Observable that emits an infinite sequence of ascending integers, with a constant interval of time of your choosing between those emissions. The first emission is not sent immediately, but only after the first period has passed.

## Example

```go
for v, _ := range rx.Interval(time.Second) {
    println(v)
}
```

Output:

```
0 // 1s
1 // 2s
2 // 3s
3 // 4s
4 // 5s
and more
```
