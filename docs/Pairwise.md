# Pairwise

> Groups pairs of consecutive emissions together and emits them as an array of two values.

## Description

> Puts the current value and previous value together as an array, and emits that.

![](https://rxjs.dev/assets/images/marble-diagrams/pairwise.png)

The Nth emission from the source Observable will cause the output Observable to emit an array [(N-1)th, Nth] of the previous and the current value, as a pair. For this reason, pairwise emits on the second and subsequent emissions from the source Observable, but not on the first emission, because there is no previous value in that case.

## Example

```go
import "fmt"

for result, err := range rxgo.Pipe1(
    rxgo.Of(1, 2, 3, 4, 5),
    rxgo.Pairwise[int](),
).Subscribe() {
    if err != nil {
        panic(err)
    }
    fmt.Println(result)
}
```

Output:

```
[1, 2]
[2, 3]
[3, 4]
[4, 5]
```
