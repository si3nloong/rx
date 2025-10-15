# Every

> Returns an Observable that emits whether or not every item of the source satisfies the condition specified.

## Description

> If all values pass predicate before the source completes, emits true before completion, otherwise emit false, then complete.

![](https://rxjs.dev/assets/images/marble-diagrams/every.png)

## Example

```go
for v, err := range rxgo.Pipe1(
    rxgo.Of(1, 1, 1, 1, 1, 3, 3),
    rxgo.Every(func(v int, _ int) bool {
        return v%2 != 0
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
true
```
