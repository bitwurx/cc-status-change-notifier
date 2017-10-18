package main

import (
	"encoding/json"
	"time"

	"github.com/satori/go.uuid"
)

// Event contains the details of a status change event.
type Event struct {
	// Kind is the type of status change event.
	// Created is the time the event occured.
	// Meta is passthrough data about the event.
	Kind    string          `json:"kind"`
	Created time.Time       `json:"created"`
	Meta    json.RawMessage `json:"meta'`
}

// NewEvent create a new event instance from the provided data.
func NewEvent(data []byte) (*Event, error) {
	evt := &Event{}
	if err := json.Unmarshal(data, evt); err != nil {
		return nil, err
	}
	return evt, nil
}

// Observer is a proxy of a remote observer instance.
type Observer struct {
	// Events is a list of all the events the observer subscribes to.
	// Id is the system derived id of the observer.
	// Conn is the socket connection for communicating with the
	// remote observer.
	Events []string `json:"events"`
	Id     string   `json:"id"`
	Conn   Conn     `json:"conn"`
}

// Notify sends the event data to the remote observer.
func (obs *Observer) Notify(evt *Event) error {
	data, _ := json.Marshal(evt)
	err := obs.Conn.WriteMessage(1, data)
	return err
}

// NewObserver creates a new observer instance.
//
// the observer is assigned a system generated version 1 uuid.
func NewObserver(data []byte, conn Conn) (*Observer, error) {
	obs := &Observer{Id: uuid.NewV1().String(), Conn: conn}
	if err := json.Unmarshal(data, obs); err != nil {
		return nil, err
	}
	return obs, nil
}
