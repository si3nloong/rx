# Max

> The `Max` operator operates on an Observable that emits numbers (or items that can be compared with a provided function), and when source Observable completes it emits a single item: the item with the largest value.

## Description

![](https://rxjs.dev/assets/images/marble-diagrams/max.png)


## Example

```go
for v, _ := range rx.Pipe1(
    rx.Of(1, 3, 4),
    rx.Max(),
) {
    println(v)
}
```

Output:

```
4
```
