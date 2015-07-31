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

This caption means for 1980 there were 19,903 children of this sex given this
name.  There were 17 names that were more popular in 1980.

![Screenshot showing separate time series for the popularity of `Jeremy` for boys and girls](images/namestats.png)

## Developing Locally

Static assets are under [static](static).  Data is stored in
[MongoDB](https://www.mongodb.org/downloads) and served from there.

To run the server locally:

```sh
go generate ./...
godep go install ./...
namestats-server --listen=:8080 --uri=mongodb://localhost/namestats
```

Now go to [http://localhost:8080](http://localhost:8080).
