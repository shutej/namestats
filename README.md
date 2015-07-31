# namestats

## Overview

Check this out at [namestats.herokuapp.com](http://namestats.herokuapp.com/).
Note that this runs on Heroku's free tier, and it's not very popular, so the
first request can take several seconds to spin up a dyno.

I wrote this to research names for my daughter.  We named her
[Isla](http://namestats.herokuapp.com/#isla).

This package project draws its name statistics from the Social Security
administration.  It shows you popularity over time.  **Hover over each bar to
see statistics for a particular year!**

```
There were 19903 births in 1980. (#18)
```

This means for that year, 19,903 children of that sex were given that name, and
there were 17 names that were more popular for that year.

![Screenshot showing separate time series for the popularity of `Jeremy` for boys and girls](images/namestats.png)

## Developing Locally

Static assets are built with [Elm](http://elm-lang.org/install) and then
minified with [uglifyjs](http://lisperator.net/uglifyjs/).  Data is served from
[MongoDB](https://www.mongodb.org/downloads).

To run the server locally:

```sh
go generate ./...
namestats-server --listen=:8080 --uri=mongodb://localhost/namestats
```

Now go to [http://localhost:8080](http://localhost:8080).
