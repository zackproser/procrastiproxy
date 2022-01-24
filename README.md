# Overview

Procrastiproxy is a proxy designed to help you focus during the day by blocking distracting websites.

It implements an in-memory, mutable list for tracking hosts that should be blocked by the proxy. This in-memory list allows for fast (`O(1)` or "constant time") look-ups.

# Running locally

`go build`

`./procrastiproxy 8001 --config .procrastiproxy.yaml`

# Admin control

Make a request to the `<server-root>/admin/` path, passing either `block` or `unblock` followed by a host, like so:

### Add a new host to the block list

`curl http://localhost:8001/admin/block/reddit.com`

### Remove a host from the block list

`curl http://localhost:8001/admin/unblock/reddit.com`

# Running tests

Procrastiproxy comes complete with tests to verify its functionality.

`go test -v ./...`
