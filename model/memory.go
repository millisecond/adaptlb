package model

import "net/http"

// Container for all state associated with an inbound request
type Request struct {
	Type string // "http", "tcp", or"udp"

	Listener *Listener
	Backend  *ServerPool

	LiveServer *LiveServer

	// If http-type
	HTTPRequest *http.Request
}

// In-memory structure to store state per-backend
type SharedState struct {
	Requests uint64
}

// In-memory structure that combines Backend and the results of Healthcheck
type LiveServer struct {
	Server *Backend

	// Healthcheck state
	Healthy             bool
	SuccessiveFailures  int
	SuccessiveSuccesses int
}
