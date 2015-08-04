Serves up a ratio of http 500's.
Ratio can be implied by the url queried or customized for a given endpoint.

Good response is a `http 200` with content `ok`.
Bad response is a `http 500` with content `panic`.

first start it:
```
$ ./bad-http-server localhost:8888
will listen for http traffic on localhost:8888
```

## fixed ratio for replies based on url

```
http://localhost:8888/reply/X
```

For 0 <= X <= 100, each reply will be such that the ratio of bad/good replies so far matches X as closely as possible.

## fixed ratio for client ip's based on url

```
http://localhost:8888/client/X
```

Creates a bucket of clients that get always good responses, and a bucket of clients that always get bad responses.
Clients seen in the last 5 minutes are partitioned into buckets based on their ip, so that
num-clients-bad/num-clients-good matches X as closely as possible.
(note for now ratio can get off balance if clients disappear and no new ones appear.)


## adjustable ratio based on url

```
http://localhost:8888/custom/KEY
http://localhost:8888/custom/KEY/
http://localhost:8888/custom/KEY/X
```

The first two forms are similar to /reply/ above, but with a default ratio of 0.
KEY must not contain slashes.
The third form updates the ratio. 0 <= X <= 100.
for example:

```
http://localhost:8888/custom/test-on-the-fly
```

Use the last configured ratio. (defaults to 0).

```
http://localhost:8888/custom/test-on-the-fly/50
```

update the `bad` ratio for `/custom/test-on-the-fly` to 50%

## status so far

```
http://localhost:8888/
```

List all used endpoints, including their path, ratio, good and bad served, as json.

