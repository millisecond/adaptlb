package listeners

import (
	"errors"
	"fmt"
	"log"
	"net"
	"strconv"
	"sync"
)

var tcpListenerMutex = &sync.Mutex{}
var tcpListeners = map[int]net.Listener{}
var tcpListenerConnections = map[int][]net.Conn{}
var tcpListenerMutexes = map[int]*sync.Mutex{}

func AddTCPPort(port int) error {
	tcpListenerMutex.Lock()
	defer tcpListenerMutex.Unlock()
	listen := ":" + strconv.Itoa(port)
	if _, pres := tcpListeners[port]; pres {
		return errors.New("Already listening on TCP " + listen)
	}
	log.Println("Opening LB TCP port", listen)
	l, err := net.Listen("tcp", listen)
	if err != nil {
		return err
	}
	tcpListeners[port] = l
	tcpListenerMutexes[port] = &sync.Mutex{}
	tcpListen(l)
	return nil
}

func RemoveTCPPort(port int) error {
	tcpListenerMutex.Lock()
	defer tcpListenerMutex.Unlock()
	if listener, pres := tcpListeners[port]; pres {
		delete(tcpListeners, port)
		err := listener.Close()
		if err != nil {
			return err
		}
		for _, c := range tcpListenerConnections[port] {
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

func tcpListen(l net.Listener) error {
	defer l.Close()
	for {
		c, err := l.Accept()
		if err != nil {
			return err
		}
		go handleTCPRequest(c)
	}
}

// Handles incoming requests.
func handleTCPRequest(c net.Conn) {
	port, err := portFromConn(c)
	if err != nil {
		log.Println("ERROR Capturing new HTTP connection:", err)
		return
	}

	func() {
		tcpListenerMutexes[port].Lock()
		defer tcpListenerMutexes[port].Unlock()
		tcpListenerConnections[port] = append(tcpListenerConnections[port], c)
	}()

	defer func() {
		tcpListenerMutexes[port].Lock()
		defer tcpListenerMutexes[port].Unlock()
		delete(tcpListenerConnections, port)
	}()

	// Make a buffer to hold incoming data.
	buf := make([]byte, 1024)
	// Read the incoming connection into the buffer.
	_, err = c.Read(buf)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}
	// Send a response back to person contacting us.
	c.Write([]byte("Message received."))
	// Close the connection when you're done with it.
	c.Close()
}
