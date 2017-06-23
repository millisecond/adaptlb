package lb

import (
	"net/http"
	"github.com/millisecond/adaptlb/model"
)

// Container for all state associated with an inbound request
type LBRequest struct {
	Type string // "http", "tcp", or"udp"

	Frontend    *model.Frontend
	SharedState *SharedState

	// The target of the load balancing
	LiveServer *LiveServer

	// If http-type
	RespontWriter http.ResponseWriter
	HTTPRequest   *http.Request
}

type Listener interface {
	Create(*model.Frontend)
	Stop()
	StopIfNot(*model.Frontend)
}

// In-memory structure to store state per-backend
type SharedState struct {
	Requests uint64
}

// In-memory structure that combines Backend and the results of Healthcheck
type LiveFrontend struct {
	FrontEnd *model.Frontend

	Listeners *[]*Listener
}

// In-memory structure that combines Backend and the results of Healthcheck
type LiveServer struct {
	Server *model.Backend

	// Healthcheck state
	Healthy             bool
	SuccessiveFailures  int
	SuccessiveSuccesses int
}
