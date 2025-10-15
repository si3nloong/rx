# From 

> Creates an Observable from an Array, an array-like object, an iterable object, or an channel.

## Description

> Converts almost anything to an Observable.

![](https://rxjs.dev/assets/images/marble-diagrams/from.png)

`From` converts various other objects and data types into Observables. It also converts a channel, an array-like, or an iterable object into an Observable that emits the items. A String, in this context, is treated as an array of characters.

## Example

```go
for v, _ := range rxgo.From([]int{1, 2, 3, 4}).Subscribe() {
    println(v)
}
```

Output:

```
1
2
3
4
```