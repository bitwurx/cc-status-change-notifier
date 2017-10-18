package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/bitwurx/jrpc2"
	"github.com/gorilla/websocket"
)

const (
	EventParseErrorCode jrpc2.ErrorCode = -32002                // event parse json rpc 2.0 error code.
	EventParseErrorMsg  jrpc2.ErrorMsg  = "Timetable not found" // event parse json rpc 2.0 error message.
)

// Conn contains methods for interacting with network sockets.
type Conn interface {
	// SetCloseHandler sets the close handler method on the conn.
	// WriteMessage sends the message to the websocket connection.
	SetCloseHandler(func(int, string) error)
	WriteMessage(int, []byte) error
}

// ApiV1 is the version 1 implementation of the rpc methods.
type ApiV1 struct {
	// notifier is the status change notifier instance.
	// upgrader is the websocket upgrader instance.
	notifier *StatusChangeNotifier
	upgrader websocket.Upgrader
}

// ObserverHandler sets up the observer connection instance.
func (api *ApiV1) ObserverHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := api.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	_, p, err := conn.ReadMessage()
	if err != nil {
		log.Println(err)
		return
	}
	if err = api.AddObserver(p, conn); err != nil {
		log.Println(err)
		conn.Close()
		return
	}
}

// AddObserver adds the observer to the api notifier.
func (api *ApiV1) AddObserver(b []byte, conn Conn) error {
	var events []string
	if err := json.Unmarshal(b, &events); err != nil {
		return err
	}
	obs, _ := NewObserver([]byte("{}"), conn)
	obs.Events = events
	api.notifier.AddObserver(obs)
	return nil
}

// NotifyParams contains the rpc parameters for the Notify method.
type NotifyParams struct {
	Kind    *string                 `json:"kind"`
	Created *time.Time              `json:"created"`
	Meta    *map[string]interface{} `json:"meta"`
}

// FromPositional parses the kind, created and meta parameters from
// the positional parameters.
func (params *NotifyParams) FromPositional(args []interface{}) error {
	if len(args) < 1 {
		return errors.New("kind parameter is required")
	} else if len(args) > 3 {
		return errors.New("kind, created, and meta parameters are required")
	}

	kind := args[0].(string)
	params.Kind = &kind

	if len(args) >= 2 && args[1] != nil {
		created, _ := time.Parse(time.RFC3339, args[1].(string))
		params.Created = &created
	}

	if len(args) == 3 && args[2] != nil {
		meta := args[2].(map[string]interface{})
		params.Meta = &meta
	}

	return nil
}

// Notify sends the incoming event to the api notifier instance.
func (api *ApiV1) Notify(params json.RawMessage) (interface{}, *jrpc2.ErrorObject) {
	e := new(Event)
	p := new(NotifyParams)
	if err := jrpc2.ParseParams(params, p); err != nil {
		return -1, err
	}
	data, _ := json.Marshal(p)
	json.Unmarshal(data, e)
	api.notifier.Notify(e)
	return 0, nil
}

// NewApiV1 creates a new instance of the version 1 api
func NewApiV1(s *jrpc2.Server) *ApiV1 {
	api := &ApiV1{
		NewStatusChangeNotifier(),
		websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin:     func(r *http.Request) bool { return true },
		},
	}
	s.Register("notify", jrpc2.Method{Method: api.Notify})
	return api
}
