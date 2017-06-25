package testutil

import (
	"fmt"
	"github.com/facebookgo/ensure"
	"github.com/millisecond/adaptlb/model"
	"log"
	"net"
	"strconv"
	"testing"
)

func TestTCPServer(t *testing.T, port int, response []byte) net.Listener {
	l, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	ensure.Nil(t, err)
	go func() {
		for {
			c, err := l.Accept()
			if err != nil {
				log.Println("Error in TCP Accept", err)
				return
			}
			go func(c net.Conn) {
				buf := make([]byte, 1024)
				// Read the incoming connection into the buffer.
				_, err = c.Read(buf)
				if err != nil {
					fmt.Println("Error reading:", err.Error())
				}
				// Send a response back to person contacting us.
				c.Write(response)
				// Close the connection when you're done with it.
				c.Close()
			}(c)
		}
	}()
	return l
}

func TCPMiniCluster(t *testing.T, responses [][]byte) []model.Backend {
	backends := []model.Backend{}
	for _, response := range responses {
		port := UniquePort()
		TestTCPServer(t, port, response)
		backends = append(backends, model.Backend{Type: "individual", Address: "localhost", Port: port})
	}
	return backends
}

func SendTCP(servAddr string, send []byte) ([]byte, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", servAddr)
	if err != nil {
		return nil, err
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return nil, err
	}

	_, err = conn.Write([]byte(send))
	if err != nil {
		return nil, err
	}

	reply := make([]byte, 1024)
	n, err := conn.Read(reply)
	if err != nil {
		return nil, err
	}

	conn.Close()
	return reply[:n], nil
}
