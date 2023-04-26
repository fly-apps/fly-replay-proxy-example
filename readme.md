# Fly Replay Proxy

An example proxy application to use on [fly.io](https://fly.io), which makes use of the [fly-replay](https://fly.io/docs/reference/dynamic-request-routing/) header
to help with dynamic request routing.

Rather than being a traditional reverse-proxy, we instead just need to return a response with a `fly-replay` header, and the Fly Proxy will do the rest.

## How to Use

This is intended for you to play with it, and adapt it if you'd like. It doesn't handle things like TLS connections nor graceful shutdowns.

### Usage

```bash
# Build for local play
go build -o bin/prox

# Create database customers.db in cwd
./bin/prox -create

# Populate the database with 5000 rows of data
./bin/prox -populate

# Run the "proxy" web server, optionally
# with a listen address (defaults to :8080)
./bin/prox -run [-addr ":8080"]
```

### Example

To play with the web server, I ran `curl` requests like this - which returns timing back to us.

The `curl-format.txt` file is included in this repository:

```bash
curl -w "@curl-format.txt" -s -i \
    -H "Host: AAqPIhi.biz" \
    http://localhost:8080/
```

An example response, with timing data output from `curl`:

```
HTTP/1.1 200 OK
Fly-Replay: instance=4955097dd2d5473982aa26c66a0c1107
Date: Wed, 26 Apr 2023 15:50:21 GMT
Content-Length: 0

     time_namelookup:  0.003899s
        time_connect:  0.004087s
     time_appconnect:  0.000000s
    time_pretransfer:  0.004103s
       time_redirect:  0.000000s
  time_starttransfer:  0.004309s
                     ----------
          time_total:  0.004364s
```

### Trying Yourself

After creating the database and generating some data, you can play with it yourself.

You can find some example (fake) Hosts ahead of time using something like this:

```bash
sqlite3 customers.db
sqlite> select * from customers order by id desc limit 10;
sqlite> ... some results here ...
sqlite> .quit
```

The `host`'s returned can be used with the Host header in your requests to get correct results from the code base via `-H "Host <somehost>"`.

## Using SQLite

This comes with an example sqlite implementation, to show how making lookups can be quite fast.

In localhost tests, a request without a SQL lookup took about 4.1ms. 

A request resulting **with** a SQL lookup against 5000 rows (and an index), took about 4.5ms. Pretty fast.

### The Database

The database `customers.db` contains a single `customers` table, creating via the following:

```sql
CREATE TABLE IF NOT EXISTS "customers" (
    "id" integer primary key autoincrement not null,
    "host" varchar not null,
    "app" varchar not null,
    "instance" varchar not null
);

CREATE UNIQUE INDEX IF NOT EXISTS "customers_host_unique" on "customers" ("host");
```
