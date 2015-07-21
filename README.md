Serves up a ratio of http 500's.
Ratio can be implied by the url queried or customized for a given endpoint.

Good response is a `http 200` with content `ok`.
Bad response is a `http 500` with content `panic`.

first start it:
```
$ ./bad-http-server localhost:8888
will listen for http traffic on localhost:8888
```

## fixed ratio based on url

```
http://localhost:8888/bad/X
```

Ratio of bad replies will always be 0 <= X <= 100.

## adjustable ratio based on url

```
http://localhost:8888/custom/KEY
http://localhost:8888/custom/KEY/
http://localhost:8888/custom/KEY/X
```

The first two forms are similar to above, but with a default ratio of 0.
KEY must not contain slashes.
The third form updates the ratio. for example:

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

