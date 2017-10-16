package main

// StatusChangeNotifier transmits subscribed events to the target
// observers.
type StatusChangeNotifier struct {
	// observers is a map of lists of observers keyed on event kind.
	observers map[string][]*Observer
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
			if nilIdx >= 0 {
				n.observers[evt][nilIdx] = obs
			} else {
				n.observers[evt] = append(n.observers[evt], obs)
			}
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
			if obs.Id == observer.Id {
				n.observers[evt][i] = nil
			}
		}
	}
}

// Notify broadcasts the event to all observers that subscribe to the
// provided event.
func (n *StatusChangeNotifier) Notify(evt *Event) {
	for _, obs := range n.observers[evt.Kind] {
		if obs != nil {
			obs.Notify(evt)
		}
	}
}

// NewStatusChangeNotifier creates a new status change notifier
// instance.
func NewStatusChangeNotifier() *StatusChangeNotifier {
	return &StatusChangeNotifier{make(map[string][]*Observer)}
}
