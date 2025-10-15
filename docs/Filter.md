# Filter

> Filter items emitted by the source Observable by only emitting those that satisfy a specified predicate.

## Description

> It only emits a value from the source if it passes a criterion function.

![](https://rxjs.dev/assets/images/marble-diagrams/filter.png)

This operator takes values from the source Observable, passes them through a predicate function and only emits those values that yielded true.

## Example

```go
for v, err := range rxgo.Pipe1(
    rxgo.Of(1, 2, 3, 4, 7, 8, 9, 10, 11, 14),
    rxgo.Filter(func(v int) bool {
        return v%2 == 0
    }),
).Subscribe() {
    if err != nil {
        panic(err)
    }
    println(v)
}
```

Output:

```
2
4
8
10
14
```
