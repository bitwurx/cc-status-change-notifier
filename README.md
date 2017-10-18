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

kind - (*String*) the event type.

created - (*String*) the creation time of the event.  
<sub><sup>*[RFC3339](https://www.ietf.org/rfc/rfc3339.txt) format (ie. 1996-12-19T16:39:57-08:00)*</sup></sub>.

runAt - (*String*) the execution point in time of the task.

#### Returns:
(*Number*) 0 on success or -1 on failure
