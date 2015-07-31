# namestats-postprocess

This command is used to postprocess name statistics gathered from the Social
Security Administration in MongoDB.

This post-processing step adds several structures that support desired use
cases, including:

* A lowercase-normalized version of the name so we can do case-insensitive
  searching.
* The total number of people with a given name, so we can rank results by
  overall popularity.
* Related names, as determined using the Levenshtein string-edit distance.

## Example Input

```json
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

## Example Output

```json
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
    },
    "lowerName": "osbaldo",
    "relatedNames": [
        "Oswaldo",
        "Osvaldo"
    ],
    "totalCount": 161
}
```

## TODO

### Write a useful usage message.

### Unit test the logic in this command.

- [ ] Write a test for `helpers.LevenshteinDistance`.
- [ ] Move logic into `../../postprocesslib`.
- [ ] Dependency inject a custom interface that supports `FindAll` and
      `UpsertId`, wrapping `mgo`.  Leave `mgo` in the `main` package.
- [ ] Use mockgen from `go generate` to create a mock for the above interface.
- [ ] Write an `postprocesslib_test.go` file that uses the mock to verify logic.
