# Concat

> Creates an output Observable which sequentially emits all values from the first given Observable and then moves on to the next.

## Description

> Concatenates multiple Observables together by sequentially emitting their values, one Observable after the other.

![](https://rxjs.dev/assets/images/marble-diagrams/concat.png)

`Concat` joins multiple Observables together, by subscribing to them one at a time and merging their results into the output Observable. You can pass either an array of Observables, or put them directly as arguments. Passing an empty array will result in Observable that completes immediately.

`Concat` will subscribe to first input Observable and emit all its values, without changing or affecting them in any way. When that Observable completes, it will subscribe to then next Observable passed and, again, emit its values. This will be repeated, until the operator runs out of Observables. When last input Observable completes, `Concat` will complete as well. At any given moment only one Observable passed to operator emits values. If you would like to emit values from passed Observables concurrently, check out merge instead, especially with optional concurrent parameter. As a matter of fact, `Concat` is an equivalent of merge operator with concurrent parameter set to 1.

Note that if some input Observable never completes, `Concat` will also never complete and Observables following the one that did not complete will never be subscribed. On the other hand, if some Observable simply completes immediately after it is subscribed, it will be invisible for `Concat`, which will just move on to the next Observable.

If any Observable in chain errors, instead of passing control to the next Observable, `Concat` will error immediately as well. Observables that would be subscribed after the one that emitted error, never will.

If you pass to `Concat` the same Observable many times, its stream of values will be "replayed" on every subscription, which means you can repeat given Observable as many times as you like. If passing the same Observable to `Concat` 1000 times becomes tedious, you can always use repeat.

## Example

```go
for v, err := range rxgo.Concat(
    rxgo.Pipe1(rxgo.Interval(time.Second), rxgo.Take(2)),
    rxgo.Pipe1(rxgo.Interval(time.Second), rxgo.Take(2)),
).Subscribe() {
    if err != nil {
        panic(err)
    }
    println(v)
}
println("Completed!")
```

Output:

```
0 // after 1s
1 // after 2s
0 // after 3s
1 // after 4s
Completed!
```
