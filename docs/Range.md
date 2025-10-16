# Range

> Creates an Observable that emits a sequence of numbers within a specified range.

## Description

> Emits a sequence of numbers in a range.

![](https://rxjs.dev/assets/images/marble-diagrams/range.png)

`Range` operator emits a range of sequential integers, in order, where you select the start of the range and its length.

## Example

```go
for v, _ := range rx.Range(1, 5).Subscribe() {
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
