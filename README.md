Serve up `ratio` of http 500's, default 0.
can be set by querying for any path that ends on an integer in the 0-100 range (inclusive)
accepts any query on any path otherwise.

example:
```
$ ./bad-http-server localhost:8888
will listen for http traffic on localhost:8888
```

```
~ curl 'http://localhost:8888/foo/whatever'
ok
~ curl 'http://localhost:8888/foo/whatever'
ok
~ curl 'http://localhost:8888/foo/whatever'
ok
~ curl 'http://localhost:8888/foo/whatever/70'
updated
~ curl 'http://localhost:8888/foo/whatever'
panic.
~ curl 'http://localhost:8888/foo/whatever'
panic.
~ curl 'http://localhost:8888/foo/whatever'
panic.
~ curl 'http://localhost:8888/foo/whatever'
ok
~ curl 'http://localhost:8888/foo/whatever'
ok
~
```
