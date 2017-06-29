# bad-http-server

Serves up a ratio of http 500's.

Good response is a `http 200` with content `ok`.
Bad response is a `http 500` with content `panic`.

Ratio can be static (implied by the url queried) or dynamic (customized over time for the same url).
The choice of which requests to affect can be made without differentiating clients, or on a per-IP basis.
So there's 2x2 = 4 different ways of working.

# Starting it
```
$ ./bad-http-server localhost:8888
will listen for http traffic on localhost:8888
```

## static ratio based on url

```
http://localhost:8888/static/X
```

For 0 <= X <= 100, each reply will be such that the ratio of bad/good replies so far matches X as closely as possible.

## static ratio based on url and client ip

```
http://localhost:8888/static-by-ip/X
```

Clients seen are partitioned into buckets based on their ip, a bucket for good responses,
one for bad responses, where num-clients-bad/num-clients-good matches X as closely as possible.

Btw, maintaining ip to bucket assignements can be tricky due to :
a) clients appearing and disappearing
b) the desire for clients to "stick" to a bucket. (a key use case of this tool is verifying [worldping](https://grafana.com/plugins/raintank-worldping-app) alerting, which is typically configured to look at streaks of errors for the same clients)
c) having a small set of clients, often only a few. this rules out consistent hashing which would otherwise be a great choice.

Buckets are rebalanced when:
a) the ratio hasn't been updated for an hour (doesn't apply here, but applies for dynamic-by-ip below)
b) clients have gone missing for at least an hour


## dynamic ratio based on url

```
http://localhost:8888/dynamic/KEY
http://localhost:8888/dynamic/KEY/
http://localhost:8888/dynamic/KEY/X
```

The first two forms are similar to /static/ above, but with a default ratio of 0.
KEY must not contain slashes.
The third form updates the ratio. 0 <= X <= 100.
for example:

```
http://localhost:8888/dynamic/test-on-the-fly
```

Use the last configured ratio. (defaults to 0).

```
http://localhost:8888/dynamic/test-on-the-fly/50
```

update the `bad` ratio for `/dynamic/test-on-the-fly` to 50%

## adjustable ratio based on url and client ip

```
http://localhost:8888/dynamic-by-ip/KEY
http://localhost:8888/dynamic-by-ip/KEY/
http://localhost:8888/dynamic-by-ip/KEY/X
```

Similar to above, customizable with default ratio 0, but bucketed by ip like static-by-ip above.

## status so far

```
http://localhost:8888/
```

List all used endpoints, including their path, ratio, good and bad served, as json.

