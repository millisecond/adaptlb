package lb

import (
	"github.com/facebookgo/ensure"
	"github.com/millisecond/adaptlb/config"
	"github.com/millisecond/adaptlb/model"
	"github.com/millisecond/adaptlb/testutil"
	"strconv"
	"testing"
)

func TestTCPListener(t *testing.T) {
	listener := &model.Frontend{Ports: "8000"}

	send := []byte("YO")
	expect := []byte("OK")

	port := 8000
	servAddr := "localhost:" + strconv.Itoa(port)
	err := AddTCPPort(listener)
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

func TestTCPActivation(t *testing.T) {
	port := "7001"
	frontend := &model.Frontend{
		Type:        "tcp",
		ServerPools: []*model.ServerPool{},
		Ports:       port,
	}

	cfg := &config.Config{
		Frontends: []*model.Frontend{frontend},
	}

	err := Activate(cfg)
	ensure.Nil(t, err)

	send := []byte("YO")
	expect := []byte("OK")

	servAddr := "localhost:" + port
	ensure.Nil(t, err)

	resp, err := testutil.SendTCP(servAddr, send)
	ensure.Nil(t, err)
	ensure.DeepEqual(t, resp, expect)

	// re-activate, make sure it's a no-op
	err = Activate(cfg)
	ensure.Nil(t, err)

	resp, err = testutil.SendTCP(servAddr, send)
	ensure.Nil(t, err)
	ensure.DeepEqual(t, resp, expect)

	// Close and make sure no more conns are accepted
	err = Activate(&config.Config{})
	ensure.Nil(t, err)

	_, err = testutil.SendTCP(servAddr, send)
	ensure.NotNil(t, err)
}
