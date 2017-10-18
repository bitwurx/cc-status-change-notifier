# Concord Status Change Notifier

Concord Status Change Notifier broadcasts resource and task status changes to observers.

### Usage
To build the docker image run:

`make build`

To run the full test suite run:

`make test`

To run the short (dependency free) test suite run:

`make test-short`

### JSON-RPC 2.0 HTTP API - Method Reference

This service uses the [JSON-RPC 2.0 Spec](http://www.jsonrpc.org/specification) over HTTP for its API.

---
#### notify(kind, created, meta) : broadcast an event to subscribers
---

#### Returns:
(*Number*) 0 on success or -1 on failure
