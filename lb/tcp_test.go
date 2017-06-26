package lb

import (
	"github.com/facebookgo/ensure"
	"github.com/millisecond/adaptlb/config"
	"github.com/millisecond/adaptlb/model"
	"github.com/millisecond/adaptlb/testutil"
	"testing"
	"time"
)

func TestTCPActivation(t *testing.T) {
	t.Parallel()

	send := []byte("YO")
	expect := []byte("OK")

	port := testutil.UniquePortString()
	backPort := testutil.UniquePort()
	testutil.TestTCPServer(t, backPort, expect)

	cfg := &config.Config{
		Frontends: []*model.Frontend{{
			RowID: "abc",
			Type:  model.LBTypeTCP,
			ServerPools: []*model.ServerPool{
				{Backends: []model.Backend{
					{Type: "individual", Address: "localhost", Port: backPort},
				}},
			},
			Ports: port,
		}},
	}

	err := Activate(nil, cfg)
	ensure.Nil(t, err)

	servAddr := "localhost:" + port
	ensure.Nil(t, err)

	resp, err := testutil.SendTCP(servAddr, send)
	ensure.Nil(t, err)
	ensure.DeepEqual(t, resp, expect)

	// Close and make sure no more conns are accepted
	err = Activate(cfg, &config.Config{})
	ensure.Nil(t, err)

	// socket operations aren't immediate
	time.Sleep(time.Millisecond * 50)

	_, err = testutil.SendTCP(servAddr, send)
	ensure.NotNil(t, err)
}

func TestTCPSingleBackend(t *testing.T) {
	t.Parallel()

	frontPort := testutil.UniquePortString()
	backPort := testutil.UniquePort()
	testutil.TestTCPServer(t, backPort, []byte("RESP"))
	cfg := &config.Config{
		Frontends: []*model.Frontend{{
			Type: model.LBTypeTCP,
			ServerPools: []*model.ServerPool{
				{Backends: []model.Backend{
					{Type: "individual", Address: "localhost", Port: backPort},
				}},
			},
			Ports: frontPort,
		}},
	}

	err := Activate(nil, cfg)
	ensure.Nil(t, err)

	send := []byte("YO")
	expect := []byte("RESP")

	servAddr := "localhost:" + frontPort
	ensure.Nil(t, err)

	resp, err := testutil.SendTCP(servAddr, send)
	ensure.Nil(t, err)
	ensure.DeepEqual(t, resp, expect)

	testutil.ValidateConsistentResponses(t, servAddr, send)
}

func TestTCPRoundRobin(t *testing.T) {
	t.Parallel()

	frontPort := testutil.UniquePortString()
	backPort := testutil.UniquePort()
	testutil.TestTCPServer(t, backPort, []byte("RESP"))
	cfg := &config.Config{
		Frontends: []*model.Frontend{{
			Type: model.LBTypeTCP,
			ServerPools: []*model.ServerPool{
				{
					Strategy: model.LBStrategyRoundRobin,
					Backends: testutil.TCPMiniCluster(t, [][]byte{[]byte("ONE"), []byte("TWO")}),
				},
			},
			Ports: frontPort,
		}},
	}

	err := Activate(nil, cfg)
	ensure.Nil(t, err)

	send := []byte("YO")

	servAddr := "localhost:" + frontPort
	ensure.Nil(t, err)

	for i := 0; i < 100; i++ {
		// LB starts at req 1, so it's the [1]th server first
		resp, err := testutil.SendTCP(servAddr, send)
		ensure.Nil(t, err)
		ensure.DeepEqual(t, resp, []byte("TWO"))

		resp, err = testutil.SendTCP(servAddr, send)
		ensure.Nil(t, err)
		ensure.DeepEqual(t, resp, []byte("ONE"))
	}
}

func TestTCPPersistence(t *testing.T) {
	t.Parallel()

	frontPort := testutil.UniquePortString()
	backPort := testutil.UniquePort()
	testutil.TestTCPServer(t, backPort, []byte("RESP"))
	cfg := &config.Config{
		Frontends: []*model.Frontend{{
			Type: model.LBTypeTCP,
			ServerPools: []*model.ServerPool{
				{
					Strategy: model.LBStrategyRoundRobin,
					Backends: testutil.TCPMiniCluster(t, [][]byte{[]byte("ONE"), []byte("TWO")}),
				},
			},
			Ports: frontPort,
		}},
	}

	err := Activate(nil, cfg)
	ensure.Nil(t, err)

	send := []byte("YO")

	servAddr := "localhost:" + frontPort
	ensure.Nil(t, err)

	for i := 0; i < 100; i++ {
		// LB starts at req 1, so it's the [1]th server first
		resp, err := testutil.SendTCP(servAddr, send)
		ensure.Nil(t, err)
		ensure.DeepEqual(t, resp, []byte("TWO"))

		resp, err = testutil.SendTCP(servAddr, send)
		ensure.Nil(t, err)
		ensure.DeepEqual(t, resp, []byte("ONE"))
	}
}
