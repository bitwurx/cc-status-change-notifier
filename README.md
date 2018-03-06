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

### Subcribing to events

A service can connect as an *observer* to the status change notifier using websockets to receive events.  The connection procedure is as follows:

- Establish a websocket connection to the status change notifier at the /observers path. *(ex. ws://status-change-notifer/observers)*
- Once the connection is established the status change notifier will accept a single message containing an array of the event types that the connecting service should be subscribed to. *(ex. ['SomeEvent', 'AnotherEvent'])*
- Once the event types array has been received the status change notifier will begin to forward events with that match the subscribed event types to the service.
- Events will be sent to the listening service with the JSON schema:

    ```
    {
        kind: "SomeEvent",
        created: "2017-01-01T12:00:00Z",
        meta: {"key1": "value1", "key2": "value2"}
    }
    ```
  The meta object will be populated with service defined data.