# Find

> Emits only the index of the first value emitted by the source Observable that meets some condition.

## Description

> It's like `Find`, but emits the index of the found value, not the value itself.

![](https://rxjs.dev/assets/images/marble-diagrams/findIndex.png)

`FindIndex` searches for the first item in the source Observable that matches the specified condition embodied by the predicate, and returns the (zero-based) index of the first occurrence in the source. Unlike `First`, the predicate is required in `FindIndex`, and does not emit an error if a valid value is not found.

## Example

```go
for v, err := range rxgo.Pipe1(
    rxgo.Of(1, 2, 3, 4, 7, 8, 9, 10, 11, 14),
    rxgo.Find(func(v int, _ int) bool {
        return v == 3
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
3
```

```go
for v, err := range rxgo.Pipe1(
    rxgo.Of(1, 2, 3, 4, 7, 8, 9, 10, 11, 14),
    rxgo.Find(func(v int, _ int) bool {
        return v > 20
    }),
).Subscribe() {
    if err != nil {
        fmt.Println(err)
        return
    }
    println(v)
}
```

Output:

```
rxgo: no values match
```
