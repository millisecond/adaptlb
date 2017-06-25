package lb

import (
	"fmt"
	"github.com/millisecond/adaptlb/model"
	"log"
	"net"
	"strconv"
	"sync"
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

//
//func RemoveTCPPort(port int) error {
//	tcpListenerMutex.Lock()
//	defer tcpListenerMutex.Unlock()
//	if listener, pres := tcpListeners[port]; pres {
//		delete(tcpListeners, port)
//		err := listener.Socket.Close()
//		if err != nil {
//			return err
//		}
//		for _, c := range listener.Connections[port] {
//			err := c.Close()
//			if err != nil {
//				return err
//			}
//		}
//		return nil
//	} else {
//		listen := ":" + strconv.Itoa(port)
//		return errors.New("Already listening on HTTP " + listen)
//	}
//}

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
func handleTCPRequest(listener *model.Listener, c net.Conn) {
	port, err := portFromConn(c)
	if err != nil {
		log.Println("ERROR Capturing new HTTP connection:", err)
		return
	}

	func() {
		listener.Mutex.Lock()
		defer listener.Mutex.Unlock()
		listener.Connections[port] = append(listener.Connections[port], c)
	}()

	defer func() {
		listener.Mutex.Lock()
		listener.Mutex.Unlock()
		delete(listener.Connections, port)
	}()

	//lbReq := &model.LBRequest{
	//	Type:     "tcp",
	//	Frontend: listener.Frontend,
	//}

	//model.LoadBalance(lbReq)

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
