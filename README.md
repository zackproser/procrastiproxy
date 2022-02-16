# Overview

Procrastiproxy is a proxy designed to help you focus during the day by blocking distracting websites.

It implements an in-memory, mutable list for tracking hosts that should be blocked by the proxy. This in-memory list allows for fast (`O(1)` or "constant time") look-ups.

# Getting started

**Use procrastiproxy as a library**

You can import procrastiproxy into your own project:

```golang
import github.com/zackproser/procrastiproxy
```

**Install procrastiproxy as a command using go**

Alternatively, you can install and use it as a command line interface (CLI) tool:

```bash
go install github.com/zackproser/procrastiproxy
```

**Install procrastiproxy as a command using the install script**

You can also run the install script to handle setup for you.

```bash
./install.sh
```

# Running locally

```bash
go build -o procrastiproxy ./cmd

./procrastiproxy --port 8000 --block reddit.com
```

You can pass multiple hosts for the `--block` flag. Just separate them with a comma like so:

```bash
./procrastiproxy --port 8000 --block reddit.com,nytimes.com
```

You can test everything is working by telling curl to use your proxy when issuing requests:

```bash
# -x tells curl to use the URL that follows as a proxy
curl -x http://localhost:8000 http://reddit.com
```

# Features

## Configurable and dynamic block list

The block list is in memory and is implemented as a map for fast lookups. You can set your baseline block list by passing the `--block` flag, like so:

```bash
procrastiproxy --port 3000 --block reddit.com,nytimes.com
```

It can be modified at runtime via the admin control endpoints described below.

## Admin control

Make a request to the `<server-root>/admin/` path, passing either `block` or `unblock` followed by a host, like so:

### Add a new host to the block list

```bash
curl http://localhost:8000/admin/block/reddit.com
```

You'll receive a response indicating success:

```bash
Successfully added: reddit.com to the block list
```

### Remove a host from the block list

```bash
curl http://localhost:8000/admin/unblock/reddit.com
```

You'll receive a response indicating success:

```bash
Successfully removed: reddit.com from the block list
```

## Office hours

You can set your working hours by passing the `--block-start-time` and `--block-end-time` flags:

```bash
procrastiproxy --block reddit.com,nytimes.com --block-start-time 9:15AM --block-end-time 5:00PM
```

If a request is made to procrastiproxy within the configured office hours, the request will be examined and blocked if its host is on the block list. If a request is made to procrastiproxy outside of the configured office hours, it will be allowed.

# Running tests

Procrastiproxy comes complete with tests.

```bash
go test -v ./...
```
