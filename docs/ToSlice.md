# ToSlice

> Collects all source emissions and emits them as an array when the source completes.

## Description

> Get all values inside an array when the source completes

![](https://rxjs.dev/assets/images/marble-diagrams/toArray.png)

`ToSlice` will wait until the source Observable completes before emitting the array containing all emissions. When the source Observable errors no array will be emitted.

## Example

```go
for v, _ := range rxgo.Pipe1(
    rxgo.Of(1, 3, 4),
    rxgo.ToSlice(),
) {
    fmt.Println(v)
}
```

Output:

```
[1, 3, 4]
```