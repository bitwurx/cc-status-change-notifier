package main

import (
	"testing"
)

func TestNewStatusChangeNotifier(t *testing.T) {
	NewStatusChangeNotifier()
}

type MockTestStatusChangeNotifierAddObserverConn struct{}

func (c MockTestStatusChangeNotifierAddObserverConn) Write(b []byte) (int, error) {
	return 0, nil
}

func TestStatusChangeNotifierAddObserver(t *testing.T) {
	conn := &MockTestStatusChangeNotifierAddObserverConn{}
	observer, err := NewObserver([]byte(`{"events": ["test"]}`), conn)
	if err != nil {
		t.Fatal(err)
	}
	notifier := NewStatusChangeNotifier()
	notifier.AddObserver(observer)
	notifier.RemoveObserver(observer)
	notifier.AddObserver(observer)
	notifier.AddObserver(observer)
	if len(notifier.observers) != 1 {
		t.Fatal("expected notifier observers count to be 1")
	}
	if notifier.observers["test"][0].Id != observer.Id {
		t.Fatal("expected notifier observer id to match observer id")
	}
}

type MockTestStatusChangeNotifierRemoveObserverConn struct{}

func (c MockTestStatusChangeNotifierRemoveObserverConn) Write(b []byte) (int, error) {
	return 0, nil
}

func TestStatusChangeNotifierRemoveObserver(t *testing.T) {
	conn := &MockTestStatusChangeNotifierAddObserverConn{}
	observer, err := NewObserver([]byte(`{"events": ["test1", "test2"]}`), conn)
	if err != nil {
		t.Fatal(err)
	}
	notifier := NewStatusChangeNotifier()
	notifier.AddObserver(observer)
	if len(notifier.observers) != 2 {
		t.Fatal("expected notifier observers count to be 1")
	}
	notifier.RemoveObserver(observer)
	if notifier.observers["test1"][0] != nil && notifier.observers["test2"][0] != nil {
		t.Fatal("expected nil values in event keys test1 and test2")
	}
}

type MockTestStatusChangeNotifierNotifyConn struct {
	Called bool
}

func (c *MockTestStatusChangeNotifierNotifyConn) Write(b []byte) (int, error) {
	c.Called = true
	return 0, nil
}

func TestStatusChangeNotifierNotify(t *testing.T) {
	conn1 := &MockTestStatusChangeNotifierNotifyConn{}
	conn2 := &MockTestStatusChangeNotifierNotifyConn{}
	observer1, err := NewObserver([]byte(`{"events": ["ev1"]}`), conn1)
	observer2, err := NewObserver([]byte(`{"events": ["ev2"]}`), conn2)
	if err != nil {
		t.Fatal(err)
	}
	notifier := NewStatusChangeNotifier()
	notifier.AddObserver(observer2)
	evt, err := NewEvent([]byte(`{"kind": "ev2"}`))
	if err != nil {
		t.Fatal(err)
	}
	notifier.Notify(evt)
	if conn1.Called == true {
		t.Fatal("expected observer1 to not have been called")
	}
	if conn2.Called == false {
		t.Fatal("expected observer2 to have been called")
	}
	notifier.AddObserver(observer1)
	notifier.RemoveObserver(observer1)
	evt, err = NewEvent([]byte(`{"kind": "ev1"}`))
	if err != nil {
		t.Fatal(err)
	}
	notifier.Notify(evt)
}
