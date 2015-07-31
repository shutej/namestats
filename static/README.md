# static

These are the static resources responsible for serving an interface over the existing API.

Static assets are built using the `go generate` command.  Logic is compiled
with [Elm](http://elm-lang.org/install) and then minified with
[uglifyjs](http://lisperator.net/uglifyjs/).

## TODO

- [ ] Eliminate the hand-rolled encoders and decoders in favor of the new
  [go2elm](https://github.com/shutej/go2elm) package.
- [ ] Fix the `options` hard-coded in [Main.elm#L16](Main.elm#L16).
- [ ] Put view logic in its own package, leaving [Main.elm](Main.elm) for
  wiring.
