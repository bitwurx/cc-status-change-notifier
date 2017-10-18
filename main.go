package main

import (
	"net/http"

	"github.com/bitwurx/jrpc2"
)

func main() {
	s := jrpc2.NewServer(":8080", "/rpc")
	api := NewApiV1(s)
	http.HandleFunc("/observers", api.ObserverHandler)
	s.Start()
}
