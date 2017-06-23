package lb

import (
	"errors"
	"fmt"
	"github.com/millisecond/adaptlb/model"
	"log"
	"net"
	"strconv"
	"sync"
)

var tcpListenerMutex = &sync.Mutex{}
var tcpListeners = map[int]*TCPListener{}

type TCPListener struct {
	port        int
	mutex       *sync.Mutex
	socket      net.Listener
	frontend    *model.Frontend
	connections map[int][]net.Conn
}

func AddTCPPort(frontend *model.Frontend) error {
	tcpListenerMutex.Lock()
	defer tcpListenerMutex.Unlock()
	ports, err := parsePorts(frontend.Ports)
	if err != nil {
		return err
	}
	for _, port := range ports {
		listen := ":" + strconv.Itoa(port)
		if _, pres := tcpListeners[port]; pres {
			//if existing
			return errors.New("Already listening on TCP " + listen)
		}
		log.Println("Opening LB TCP port", listen)
		socket, err := net.Listen("tcp", listen)
		if err != nil {
			return err
		}
		tcpListener := &TCPListener{
			port:        port,
			mutex:       &sync.Mutex{},
			socket:      socket,
			frontend:    frontend,
			connections: map[int][]net.Conn{},
		}
		tcpListeners[port] = tcpListener
		go tcpListen(tcpListener)
	}
	return nil
}

func RemoveTCPPort(port int) error {
	tcpListenerMutex.Lock()
	defer tcpListenerMutex.Unlock()
	if listener, pres := tcpListeners[port]; pres {
		delete(tcpListeners, port)
		err := listener.socket.Close()
		if err != nil {
			return err
		}
		for _, c := range listener.connections[port] {
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

func tcpListen(listener *TCPListener) error {
	defer listener.socket.Close()
	for {
		c, err := listener.socket.Accept()
		if err != nil {
			return err
		}
		go handleTCPRequest(listener, c)
	}
}

// Handles incoming requests.
func handleTCPRequest(listener *TCPListener, c net.Conn) {
	port, err := portFromConn(c)
	if err != nil {
		log.Println("ERROR Capturing new HTTP connection:", err)
		return
	}

	func() {
		listener.mutex.Lock()
		defer listener.mutex.Unlock()
		listener.connections[port] = append(listener.connections[port], c)
	}()

	defer func() {
		listener.mutex.Lock()
		listener.mutex.Unlock()
		delete(listener.connections, port)
	}()

	lbReq := &model.LBRequest{
		Type:     "tcp",
		Frontend: listener.frontend,
	}

	model.LoadBalance(lbReq)

	// Make a buffer to hold incoming data.
	buf := make([]byte, 1024)
	// Read the incoming connection into the buffer.
	_, err = c.Read(buf)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}

	// Send a response back to person contacting us.
	c.Write([]byte("OK"))
	// Close the connection when you're done with it.
	c.Close()
}
