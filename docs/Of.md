# Of

> Converts the arguments to an observable sequence.

## Description

> Each argument becomes a next notification.

![](https://rxjs.dev/assets/images/marble-diagrams/of.png)

Unlike `From`, it does not do any flattening and emits each argument in whole as a separate `next` notification.

## Example

```go
for v, _ := range rx.Of(1, 2, 3, 4, 5) {
    println(v)
}
```

Output:

```
1
2
3
4
5
```
