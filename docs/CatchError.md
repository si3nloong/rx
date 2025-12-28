# CatchError 

> Catches errors on the observable to be handled by returning a new observable or throwing an error.

## Description

> It only listens to the error channel and ignores notifications. Handles errors from the source observable, and maps them to a new observable. The error may also be rethrown, or a new error can be thrown to emit an error from the result.

![](https://rxjs.dev/assets/images/marble-diagrams/catch.png)

This operator handles errors, but forwards along all other events to the resulting observable. If the source observable terminates with an error, it will map that error to a new observable, subscribe to it, and forward all of its events to the resulting observable.


## Example

```go
for v, err := range rx.Pipe2(
    rx.Of(1, 2, 3, 4, 5),
    rx.MapErr(func(v int, _ int) (int, error) {
        if v == 4 {
            return 0, errors.New(`four`)
        }
        return v, nil
    }),
    rx.CatchError2[int](func(err error) rx.Observable[string] {
        return rx.From[string]([]string{"I", "II", "III", "IV", "V"})
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
1
2
3
I
II
III
IV
V
```