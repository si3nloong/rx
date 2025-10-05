# Max

> The `Max` operator operates on an Observable that emits numbers (or items that can be compared with a provided function), and when source Observable completes it emits a single item: the item with the largest value.

## Description

![](https://rxjs.dev/assets/images/marble-diagrams/max.png)


## Example

```go
for v, _ := range rxgo.Pipe1(
    rxgo.Of(1, 3, 4),
    rxgo.Max(),
) {
    println(v)
}
```

Output:

```
4
```
