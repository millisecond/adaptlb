package lb

import (
	"testing"
	"github.com/facebookgo/ensure"
	"github.com/millisecond/adaptlb/testutil"
	"github.com/millisecond/adaptlb/config"
	"github.com/millisecond/adaptlb/model"
)

func BenchmarkTCPConnections(b *testing.B) {
	frontPort := testutil.UniquePortString()
	backPort := testutil.UniquePort()
	testutil.TestTCPServer(b, backPort, []byte("RESP"))
	cfg := &config.Config{
		Frontends: []*model.Frontend{{
			Type: model.LBTypeTCP,
			ServerPools: []*model.ServerPool{
				{
					Strategy: model.LBStrategyRoundRobin,
					Backends: testutil.TCPMiniCluster(b, [][]byte{[]byte("ONE"), []byte("TWO")}),
				},
			},
			Ports: frontPort,
		}},
	}

	err := Activate(nil, cfg)
	ensure.Nil(b, err)

	send := []byte("YO")

	servAddr := "localhost:" + frontPort
	ensure.Nil(b, err)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// LB starts at req 1, so it's the [1]th server first
		testutil.SendTCP(servAddr, send)
	}
}
