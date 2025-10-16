# Map

> Applies a given project function to each value emitted by the source Observable, and emits the resulting values as an Observable.

## Description

> It passes each source value through a transformation function to get corresponding output values.

![](https://rxjs.dev/assets/images/marble-diagrams/map.png)

This operator applies a projection to each value and emits that projection in the output Observable.

## Example

```go
for v, err := range rx.Pipe1(
    rx.Of(1, 88, 100),
    rx.Map(func(v int, _ int) string {
        return strconv.Itoa(v)
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
"1"
"88"
"100"
```
