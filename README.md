# namestats

This serves baby name statistics from the Social Security administration.

To use this package:

```sh
go generate ./...
namestats-server --listen=:8080 --uri=mongodb://localhost/namestats
```

Now go to [http://localhost:8080](http://localhost:8080).
