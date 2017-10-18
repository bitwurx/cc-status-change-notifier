package main

import (
	"net/http"
	"testing"

	"github.com/bitwurx/jrpc2"
)

type MockApiV1AddObserverConn struct {
	Called *bool
}

func (c MockApiV1AddObserverConn) SetCloseHandler(func(int, string) error) {}

func (c MockApiV1AddObserverConn) WriteMessage(t int, p []byte) error {
	*c.Called = true
	return nil
}

func TestUpgraderCheckOrigin(t *testing.T) {
	api := NewApiV1(jrpc2.NewServer("", ""))
	if api.upgrader.CheckOrigin(&http.Request{}) != true {
		t.Fatal("expected origin check to return true")
	}
}

func TestApiV1Notify(t *testing.T) {
	called := false
	conn := &MockApiV1AddObserverConn{&called}
	api := NewApiV1(jrpc2.NewServer("", ""))
	api.AddObserver([]byte(`["jobStatusChange"]`), conn)
	result, errObj := api.Notify([]byte(`[]`))
	if errObj == nil || errObj.Code != jrpc2.InvalidParamsCode {
		t.Fatal("expected invalid params error")
	}
	result, errObj = api.Notify([]byte(`["test", {}, "uh-oh", 1]`))
	if errObj == nil || errObj.Code != jrpc2.InvalidParamsCode {
		t.Fatal("expected invalid params error")
	}
	result, errObj = api.Notify([]byte(`["test", null]`))
	if errObj != nil {
		t.Fatal(errObj.Message)
	}
	result, errObj = api.Notify([]byte(`{"test": "jobStatusChange", "created": "2017-01-01T12:00:00Z" "this": {"status": 1}}`))
	if errObj == nil || errObj.Code != jrpc2.InvalidParamsCode {
		t.Fatal("expected invalid params error")
	}
	result, errObj = api.Notify([]byte(`["jobStatusChange", "2017-01-01T12:00:00Z", {"status": 1}]`))
	if errObj != nil {
		t.Fatal(errObj.Message)
	}
	if result != 0 {
		t.Fatal("expected result to be 0")
	}
	if called != true {
		t.Fatal("expected observer to have been called")
	}
}

func TestApiV1AddObserver(t *testing.T) {
	api := NewApiV1(jrpc2.NewServer("", ""))
	if err := api.AddObserver([]byte(`"test"`), MockApiV1AddObserverConn{}); err == nil {
		t.Fatal("expected JSON unmarshal error")
	}
	if err := api.AddObserver([]byte(`["ev1"]`), MockApiV1AddObserverConn{}); err != nil {
		t.Fatal(err)
	}
	if len(api.notifier.observers["ev1"]) != 1 {
		t.Fatal("expected observers count to be 1")
	}
}
