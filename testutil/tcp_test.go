package testutil

import (
	"github.com/facebookgo/ensure"
	"testing"
)

func TestStoppableTCPListener(t *testing.T) {
	l := TestTCPServer(t,8002)
	send := []byte("YO")
	expect := []byte("OK")
	servAddr := "localhost:8002"

	resp, err := SendTCP(servAddr, send)
	ensure.Nil(t, err)
	ensure.DeepEqual(t, resp, expect)

	// Close and make sure no more conns are accepted
	err = l.Close()
	ensure.Nil(t, err)

	_, err = SendTCP(servAddr, send)
	ensure.NotNil(t, err)
}
