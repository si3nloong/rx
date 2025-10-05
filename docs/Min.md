# Min

> The `Min` operator operates on an Observable that emits numbers (or items that can be compared with a provided function), and when source Observable completes it emits a single item: the item with the smallest value.

## Description

![](https://rxjs.dev/assets/images/marble-diagrams/min.png)


## Example

```go
for v, _ := range rxgo.Pipe1(
    rxgo.Of(-88, 1, 3, 4),
    rxgo.Max(),
) {
    println(v)
}
```

Output:

```
-88
```
