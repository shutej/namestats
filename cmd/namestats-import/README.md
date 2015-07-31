# namestats-import

This command is used to import name statistics from the Social Security
Administration into MongoDB.

Input is presented as multiple pages, one per year, formatted as an HTML
table with male and female names descending in popularity.

Output is written to a MongoDB database as individual records stored in the
same collection with gender information in the record.  Data is "transposed" to
be keyed by name first and date second.  We use a Series to store a compact
sequence of integers and only one yearly offset, this minimizes payload sizes
between the database, the HTTP server, and the Elm client.

The strategy taken by this tool is "slow" (taking minutes) but idempotent
and reliable.  This tool is run from a laptop, once per year, so that's the
right trade-off.

## Example Input

```
                         Popularity in 2013
Rank    Male name    Percent of total males    Female name    Percent of total females
1       Noah         0.9040%                   Sophia         1.1022%
                                 ...
10      Daniel       0.7071%                   Elizabeth      0.4906%
```

## Example Output

```javascript
{
    "_id": {
        "$oid": "558edee34238552d85001c05"
    },
    "gender": "m",
    "name": "Osbaldo",
    "rank": {
        "key": 2002,
        "data": [
            980
        ]
    },
    "count": {
        "key": 2002,
        "data": [
            161
        ]
    }
}
```

## TODO

### Write a useful usage message.

### Unit test the logic in this command.

- [ ] Move logic into `../../importlib`.
- [ ] Insert an `io.Reader` seam between the `http` client and the parser.
      Leave client in the `main` package.
- [ ] Dependency inject a custom interface that supports `FindOne` and
      `UpsertId`, wrapping `mgo`.  Leave `mgo` in the `main` package.
- [ ] Use mockgen from `go generate` to create a mock for the above interface.
- [ ] Write an `importlib_test.go` file that uses the input and output seam to verify logic.
