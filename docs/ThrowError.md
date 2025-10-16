# ThrowError

> Converts the arguments to an observable sequence.

## Description

> Just errors and does nothing else

![](https://rxjs.dev/assets/images/marble-diagrams/throw.png)

This creation function is useful for creating an observable that will create an error and error every time it is subscribed to. Generally, inside of most operators when you might want to return an errored observable, this is unnecessary.

## Example

```go
for v, err := range rx.ThrowError(func() error {
    return errors.New(`stop`)
}).Subscribe() {
    if err != nil {
        panic("has error")
    }
    println(v)
}
```

Output:

```
panic: has error
```
