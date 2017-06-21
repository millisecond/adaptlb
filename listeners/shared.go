package listeners

import (
	"strconv"
	"net"
)

func portFromConn(c net.Conn) (int, error) {
	_, portS, err := net.SplitHostPort(c.LocalAddr().String())
	if err != nil {
		return -1, err
	}
	port, err := strconv.Atoi(portS)
	if err != nil {
		return -1, err
	}
	return port, nil
}
