package model

import (
	"log"
	"net"
	"net/http"
	"sync"
	"sync/atomic"
)

// Container for all state associated with an inbound request
type LBRequest struct {
	// Inbound reqs should assign these fields and then call LoadBalance()
	Type     string // "http", "tcp", or"udp"
	Frontend *Frontend

	// The target of the load balancing
	ServerPool *ServerPool
	LiveServer *LiveServer

	// If http-type
	ResponseWriter http.ResponseWriter
	HTTPRequest   *http.Request
}

type Listener struct {
	Secure      bool
	Port        int
	Mutex       *sync.Mutex
	Socket      net.Listener
	Frontend    *Frontend
	Connections map[int][]net.Conn
	//Create(*Frontend)
	//Stop()
	//StopIfNot(*Frontend)
	//Connections()[]*LiveConnection
}

func (listener *Listener) Stop() {
	err := listener.Socket.Close()
	if err != nil {
		log.Println(err)
	}
	for _, conns := range listener.Connections {
		for _, conn := range conns {
			err := conn.Close()
			if err != nil {
				log.Println(err)
			}
		}
	}
}

// In-memory structure to store state per-backend
type SharedLBState struct {
	Requests uint64
}

func (ss *SharedLBState) IncrAndGetRequests() uint64 {
	return atomic.AddUint64(&ss.Requests, 1)
}

// In-memory structure that combines Backend and the results of Healthcheck
type LiveServer struct {
	Address string // computed ip:port value

	// Healthcheck state
	Healthy             bool
	SuccessiveFailures  int
	SuccessiveSuccesses int
}

type LiveConnection struct {
	Conn net.Conn
}
