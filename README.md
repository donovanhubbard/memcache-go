# memcache-go

A memcached cli written in go.

## Running

Specify the host and port with either the `-host` and `-port`
switches or use the `HOST` and `PORT` environmnet variables.
Defaults to localhost port 11211.

### Retrieve a key

```bash
bin/memcache-go -command get -key foo
bar
```

### Set a key

```bash
bin/memcache-go -command set -key foo -expiry 500 -flags 0 -value bar
```

### Delete a key

```bash
bin/memcache-go -command delete -key foo
```

## Testing

The integrations test requires a working instance of memcached
on localhost port 11211.
