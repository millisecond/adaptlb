package lb

import (
	"context"
	"github.com/millisecond/adaptlb/model"
	"github.com/millisecond/adaptlb/util"
	"log"
	"net"
	"strconv"
)

var tcpListenerMutex = &util.WrappedMutex{}

func addTCPPort(frontend *model.Frontend) error {
	ctx := context.Background()
	tcpListenerMutex.Lock(ctx)
	defer tcpListenerMutex.Unlock(ctx)
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
			Mutex:       &util.WrappedMutex{},
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

	ctx := context.Background()
	func() {
		listener.Mutex.Lock(ctx)
		listener.Connections[port] = append(listener.Connections[port], inboundConn)
		listener.Mutex.Unlock(ctx)
	}()

	defer func() {
		listener.Mutex.Lock(ctx)
		delete(listener.Connections, port)
		listener.Mutex.Unlock(ctx)
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

	backendConn, err := net.Dial("tcp", lbReq.LiveServer.Address)

	// write to dst what it reads from src
	var cp = func(src, dst net.Conn) {
		defer func() {
			inboundConn.Close()
			backendConn.Close()
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
	go cp(inboundConn, backendConn)
	go cp(backendConn, inboundConn)
}
