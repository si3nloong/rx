# BufferCount 

> Buffers the source Observable values until the size hits the maximum bufferSize given.

## Description

> Collects values from the past as an array, and emits that array only when its size reaches `bufferSize`.

![](https://rxjs.dev/assets/images/marble-diagrams/bufferCount.png)

Buffers a number of values from the source Observable by bufferSize then emits the buffer and clears it, and starts a new buffer each startBufferEvery values. If startBufferEvery is not provided or is null, then new buffers are started immediately at the start of the source and when each buffer closes and is emitted.

## Example

```go
import "fmt"

for v, err := range rxgo.Pipe1(
    rxgo.From[int]([]int{1, 3, 4, 5, 9}),
    rxgo.BufferCount[int](2),
).Subscribe() {
    if err != nil {
        panic(err)
    }
    fmt.Println(v)
}
```

Output:

```
[1, 3]
[4, 5]
[9]
```