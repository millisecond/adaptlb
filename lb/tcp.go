package lb

import (
	"github.com/millisecond/adaptlb/model"
	"log"
	"net"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

var tcpListenerMutex = &sync.Mutex{}

func addTCPPort(frontend *model.Frontend) error {
	tcpListenerMutex.Lock()
	defer tcpListenerMutex.Unlock()
	ports, err := parsePorts(frontend.Ports)
	if err != nil {
		return err
	}
	for _, port := range ports {
		listen := ":" + strconv.Itoa(port)
		log.Println("Opening LB TCP port", listen)
		socket, err := net.Listen("tcp", listen)
		if err != nil {
			return err
		}
		listener := &model.Listener{
			Port:        port,
			Mutex:       &sync.Mutex{},
			Socket:      socket,
			Frontend:    frontend,
			Connections: map[int][]net.Conn{},
		}
		(*frontend.Listeners)[port] = listener
		go tcpListen(listener)
	}
	return nil
}

func tcpListen(listener *model.Listener) error {
	defer listener.Socket.Close()
	for {
		c, err := listener.Socket.Accept()
		if err != nil {
			return err
		}
		go handleTCPRequest(listener, c)
	}
}

// Handles incoming requests.
func handleTCPRequest(listener *model.Listener, inboundConn net.Conn) {
	port, err := portFromConn(inboundConn)
	if err != nil {
		log.Println("ERROR Capturing new HTTP connection:", err)
		return
	}

	func() {
		listener.Mutex.Lock()
		defer listener.Mutex.Unlock()
		listener.Connections[port] = append(listener.Connections[port], inboundConn)
	}()

	defer func() {
		listener.Mutex.Lock()
		listener.Mutex.Unlock()
		delete(listener.Connections, port)
	}()

	lbReq := &model.LBRequest{
		Type:     "tcp",
		Frontend: listener.Frontend,
	}

	validTarget := LoadBalanceL4(lbReq)
	if !validTarget {
		// When doing L4, not much we can do if we don't have targets
		inboundConn.Close()
		return
	}

	// pipeDone counts closed pipe
	var pipeDone int32
	var timer *time.Timer

	backendConn, err := net.Dial("tcp", lbReq.LiveServer.Address)

	// write to dst what it reads from src
	var pipe = func(src, dst net.Conn) {
		defer func() {
			// if it is the first pipe to end...
			if v := atomic.AddInt32(&pipeDone, 1); v == 1 {
				// ...wait 'timeout' seconds before closing connections
				timer = time.AfterFunc(time.Second, func() {
					// test if the other pipe is still alive before closing conn
					if atomic.AddInt32(&pipeDone, 1) == 2 {
						inboundConn.Close()
						backendConn.Close()
					}
				})
			} else if v == 2 {
				inboundConn.Close()
				backendConn.Close()
				timer.Stop()
			}
		}()

		buff := make([]byte, 65535)
		for {
			n, err := src.Read(buff)
			if err != nil {
				return
			}
			b := buff[:n]
			n, err = dst.Write(b)
			if err != nil {
				return
			}
		}
	}
	go pipe(inboundConn, backendConn)
	go pipe(backendConn, inboundConn)
}
