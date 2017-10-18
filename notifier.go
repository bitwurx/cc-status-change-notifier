package main

import (
	"fmt"
	"log"
)

// StatusChangeNotifier transmits subscribed events to the target
// observers.
type StatusChangeNotifier struct {
	// observers is a map of lists of observers keyed on event kind.
	observers map[string][]*Observer
}

// CloseHandler is called the connection is closed by the client.
func (n *StatusChangeNotifier) CloseHandler(obs *Observer) func(int, string) error {
	return func(code int, text string) error {
		log.Println(code, text)
		n.RemoveObserver(obs)
		return nil
	}
}

// AddObserver adds the observer to each event they subscribe to.
//
// This method populates nil indexes left over from disconnected
// observers.
func (n *StatusChangeNotifier) AddObserver(obs *Observer) {
	for _, evt := range obs.Events {
		exists := false
		nilIdx := -1
		for i, observer := range n.observers[evt] {
			if observer == nil {
				nilIdx = i
			} else if obs.Id == observer.Id {
				exists = true
			}
		}
		if !exists {
			obs.Conn.SetCloseHandler(n.CloseHandler(obs))
			if nilIdx >= 0 {
				n.observers[evt][nilIdx] = obs
			} else {
				n.observers[evt] = append(n.observers[evt], obs)
			}
			log.Println("added observer:", obs)
		}
	}
}

// RemoveObserver removes the observer from all subscribed event keys.
//
// This method leaves a nil value in the index where the observer
// was once assigned.
func (n *StatusChangeNotifier) RemoveObserver(obs *Observer) {
	for _, evt := range obs.Events {
		for i, observer := range n.observers[evt] {
			if observer != nil && obs.Id == observer.Id {
				n.observers[evt][i] = nil
				log.Println("removed observer:", obs.Id)
			}
		}
	}
}

// Notify broadcasts the event to all observers that subscribe to the
// provided event.
func (n *StatusChangeNotifier) Notify(evt *Event) {
	log.Println(fmt.Sprintf("received event: [%v] %v - %v",
		evt.Kind, evt.Created, string(evt.Meta)))

	for _, obs := range n.observers[evt.Kind] {
		if obs != nil {
			if err := obs.Notify(evt); err != nil {
				log.Println("notify error:", err)
				n.RemoveObserver(obs)
			} else {
				log.Println(fmt.Sprintf("notify [%v] - event %v", obs.Id, evt))
			}
		}
	}
}

// NewStatusChangeNotifier creates a new status change notifier
// instance.
func NewStatusChangeNotifier() *StatusChangeNotifier {
	return &StatusChangeNotifier{make(map[string][]*Observer)}
}
