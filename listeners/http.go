package listeners

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"sync"
)

var httpListenerMutex = &sync.Mutex{}
var httpListeners = map[int]net.Listener{}
var httpListenerConnections = map[int][]net.Conn{}
var httpListenerMutexes = map[int]*sync.Mutex{}

// HTTP Listeners are a single http.Server listening on multiple connections.
var HTTPServer = &http.Server{
	Handler:   http.HandlerFunc(hostnameMultiplexer),
	ConnState: connState,
}

// Requests first come into the HostnameMultiplexer which will assign it to a specific Listener (or Bad Gateway)
func hostnameMultiplexer(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "OK")
}

func connState(c net.Conn, state http.ConnState) {
	if state == http.StateNew || state == http.StateClosed {
		port, err := portFromConn(c)
		if err != nil {
			log.Println("ERROR Capturing new connection:", err)
			return
		}
		httpListenerMutexes[port].Lock()
		defer httpListenerMutexes[port].Unlock()
		if state == http.StateNew {
			httpListenerConnections[port] = append(httpListenerConnections[port], c)
		} else if state == http.StateClosed {
			delete(httpListenerConnections, port)
		}
	}
}

func AddHTTPPort(port int) error {
	httpListenerMutex.Lock()
	defer httpListenerMutex.Unlock()
	listen := ":" + strconv.Itoa(port)
	if _, pres := httpListeners[port]; pres {
		return errors.New("Already listening on HTTP " + listen)
	}
	log.Println("Opening LB HTTP port", listen)
	l, err := net.Listen("tcp", listen)
	if err != nil {
		return err
	}
	httpListeners[port] = l
	httpListenerMutexes[port] = &sync.Mutex{}
	go HTTPServer.Serve(l)
	return nil
}

func RemoveHTTPPort(port int) error {
	httpListenerMutex.Lock()
	defer httpListenerMutex.Unlock()
	if listener, pres := httpListeners[port]; pres {
		delete(httpListeners, port)
		err := listener.Close()
		if err != nil {
			return err
		}
		for _, c := range httpListenerConnections[port] {
			err := c.Close()
			if err != nil {
				return err
			}
		}
		return nil
	} else {
		listen := ":" + strconv.Itoa(port)
		return errors.New("Already listening on HTTP " + listen)
	}
}