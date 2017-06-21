package testutil

import (
	"net"
	"strconv"
	"testing"
	"github.com/facebookgo/ensure"
	"fmt"
	"log"
)

func TestTCPServer(t *testing.T, port int) net.Listener {
	l, err := net.Listen("tcp", ":" + strconv.Itoa(port))
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
				c.Write([]byte("OK"))
				// Close the connection when you're done with it.
				c.Close()
			}(c)
		}
	}()
	return l
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
