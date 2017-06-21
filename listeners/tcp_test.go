package listeners

import (
	"github.com/facebookgo/ensure"
	"strconv"
	"testing"
	"github.com/millisecond/adaptlb/model"
	"github.com/millisecond/adaptlb/testutil"
)

func TestTCPListener(t *testing.T) {
	listener := &model.Listener{}

	send := []byte("YO")
	expect := []byte("OK")

	port := 8000
	servAddr := "localhost:"+strconv.Itoa(port)
	err := AddTCPPort(listener, port)
	ensure.Nil(t, err)

	resp, err := testutil.SendTCP(servAddr, send)
	ensure.Nil(t, err)
	ensure.DeepEqual(t, resp, expect)

	// Close and make sure no more conns are accepted
	err = RemoveTCPPort(port)
	ensure.Nil(t, err)

	_, err = testutil.SendTCP(servAddr, send)
	ensure.NotNil(t, err)
}

func TestTCPLB(t *testing.T) {
	listener := &model.Listener{}

	send := []byte("YO")
	expect := []byte("OK")

	port := 8000
	servAddr := "localhost:"+strconv.Itoa(port)
	err := AddTCPPort(listener, port)
	ensure.Nil(t, err)

	resp, err := testutil.SendTCP(servAddr, send)
	ensure.Nil(t, err)
	ensure.DeepEqual(t, resp, expect)

	// Close and make sure no more conns are accepted
	err = RemoveTCPPort(port)
	ensure.Nil(t, err)

	_, err = testutil.SendTCP(servAddr, send)
	ensure.NotNil(t, err)
}
