package model

import (
	"net/http"
	"sync"
	"net"
)

// Container for all state associated with an inbound request
type LBRequest struct {
	Type string // "http", "tcp", or"udp"

	Frontend    *Frontend
	SharedState *SharedState

	// The target of the load balancing
	LiveServer *LiveServer

	// If http-type
	RespontWriter http.ResponseWriter
	HTTPRequest   *http.Request
}

type Listener interface {
	Create(*Frontend)
	Stop()
	StopIfNot(*Frontend)
}

// In-memory structure to store state per-backend
type SharedState struct {
	Requests uint64
}

// In-memory structure that combines Backend and the results of Healthcheck
type LiveFrontend struct {
	FrontEnd *Frontend

	Listeners *[]*Listener
}

// In-memory structure that combines Backend and the results of Healthcheck
type LiveServer struct {
	Server *Backend

	// Healthcheck state
	Healthy             bool
	SuccessiveFailures  int
	SuccessiveSuccesses int
}

type TCPListener struct {
	Port        int
	Mutex       *sync.Mutex
	Socket      net.Listener
	Frontend    *Frontend
	Connections map[int][]net.Conn
}
