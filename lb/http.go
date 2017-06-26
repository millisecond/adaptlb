package lb

import (
	"context"
	"errors"
	"fmt"
	"github.com/millisecond/adaptlb/model"
	"github.com/millisecond/adaptlb/util"
	"log"
	"net"
	"net/http"
	"strconv"
)

var httpListenerMutex = &util.WrappedMutex{}
var httpListeners = map[int]net.Listener{}
var httpListenerConnections = map[int][]net.Conn{}
var httpListenerMutexes = map[int]*util.WrappedMutex{}

// HTTP Listeners are a single http.Server listening on multiple connections.
var HTTPServer = &http.Server{
	Handler:   http.HandlerFunc(hostnameMultiplexer),
	ConnState: connState,
}

// Requests first come into the HostnameMultiplexer which will assign it to a specific Frontend (or Bad Gateway)
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
		ctx := context.Background()
		httpListenerMutexes[port].Lock(ctx)
		defer httpListenerMutexes[port].Unlock(ctx)
		if state == http.StateNew {
			httpListenerConnections[port] = append(httpListenerConnections[port], c)
		} else if state == http.StateClosed {
			delete(httpListenerConnections, port)
		}
	}
}

func AddHTTPPort(frontend *model.Frontend) error {
	ctx := context.Background()
	httpListenerMutex.Lock(ctx)
	defer httpListenerMutex.Unlock(ctx)
	ports, err := parsePorts(frontend.Ports)
	if err != nil {
		return err
	}
	for _, port := range ports {
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
		httpListenerMutexes[port] = &util.WrappedMutex{}
		go HTTPServer.Serve(l)
	}
	return nil
}

func RemoveHTTPPort(port int) error {
	ctx := context.Background()
	httpListenerMutex.Lock(ctx)
	defer httpListenerMutex.Unlock(ctx)
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
