# Race

> Returns an observable that mirrors the first source observable to emit an item.

## Description

![](https://rxjs.dev/assets/images/marble-diagrams/race.png)

`Race` returns an observable, that when subscribed to, subscribes to all source observables immediately. As soon as one of the source observables emits a value, the result unsubscribes from the other sources. The resulting observable will forward all notifications, including error and completion, from the "winning" source observable.

If one of the used source observable throws an errors before a first notification the `Race` operator will also throw an error, no matter if another source observable could potentially win the `Race`.

`Race` can be useful for selecting the response from the fastest network connection for HTTP or WebSockets. `Race` can also be useful for switching observable context based on user input.

## Example

```go
for result, err := range rxgo.Race(
    rxgo.Pipe1(rxgo.Interval(time.Second*7), rxgo.Map(func(_ int, _ int) string { return "slow one" })),
    rxgo.Pipe1(rxgo.Interval(time.Second*3), rxgo.Map(func(_ int, _ int) string { return "fast one" })),
    rxgo.Pipe1(rxgo.Interval(time.Second*5), rxgo.Map(func(_ int, _ int) string { return "medium one" })),
).Subscribe() {
    if err != nil {
        panic(err)
    }
    fmt.Println(result)
}
```

Output:

```
a series of 'fast one'
```
