# namestats-server

This is the main server used to serve both API requests and static assets for
the project.

* `/` serves static assets that have been compiled with `go generate`.  See the
  [static](../../static) directory for static assets.
* `/v1/namesearch/` allows us to see those names that are most popular, without running a search.
* `/v1/namesearch/:query` performs a case-insensitive prefix query for name statistics.

## TODO

### Write a useful usage message.
