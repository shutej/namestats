# babynames

This serves baby name statistics from the Social Security administration.

To use this package:

```sh
go generate ./...
babynames-server --listen=:8080
```

Now go to [http://localhost:8080](http://localhost:8080).

## TODO

* Move almost all of the Go code out of `static`, leaving `main.js`.  This is
  so the server won't statically serve `*.go`.
* Get the global max number of people ever given any name, use that for all
  bargraphs so they're on the same axis.
* Use GORM with a LIKE operator, the existing data structure sucks.
* Figure out ViewModel RPC subscriptions.
  * I think these are ViewModels that act like
    [service/keepalive/client](service/keepalive/client), subscribing to a
    stream of updates for an object...
  * Use `addDisposeCallback` to call `call.CloseStream` to clean up on the server side.
  * The call's arguments should be unpacked from the `Params` object.
