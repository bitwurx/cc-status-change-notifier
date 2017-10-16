package main

import (
	"testing"

	"github.com/satori/go.uuid"
)

type MockNotifyConn struct {
	Called bool
}

func (c *MockNotifyConn) Write(b []byte) (int, error) {
	c.Called = true
	return 0, nil
}

func TestNewEvent(t *testing.T) {
	_, err := NewEvent([]byte(`{"kind": "test"`))
	if err == nil {
		t.Fatal("expected json parse error")
	}
	_, err = NewEvent([]byte(`{"kind": "test"}`))
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewObserver(t *testing.T) {
	_, err := NewObserver([]byte(`"events": ["test",`), nil)
	if err == nil {
		t.Fatal("expected json parse error")
	}
	observer, err := NewObserver([]byte(`{"events": ["test"]}`), nil)
	if err != nil {
		t.Fatal(err)
	}
	id, err := uuid.FromString(observer.Id)
	if err != nil {
		t.Fatal(err)
	}
	if id.Version() != 1 {
		t.Fatal("expected observer id to be uuid version 1")
	}
}

func TestObserverNotify(t *testing.T) {
	conn := &MockNotifyConn{}
	evt, err := NewEvent([]byte(`{"kind": "test", "created": "2017-01-01T00:00:00Z", "meta": {}}`))
	if err != nil {
		t.Fatal(err)
	}
	observer, err := NewObserver([]byte(`{"events": ["test"], "id": "123"}`), conn)
	if err != nil {
		t.Fatal(err)
	}
	if err = observer.Notify(evt); err != nil {
		t.Fatal(err)
	}
	if conn.Called != true {
		t.Fatal("expected observer Write to have been called")
	}
}
