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

	port := testutil.UniquePortString()
	frontend := &model.Frontend{
		RowID:       "abc",
		Type:        "tcp",
		ServerPools: []*model.ServerPool{},
		Ports:       port,
	}

	cfg := &config.Config{
		Frontends: []*model.Frontend{frontend},
	}

	err := Activate(nil, cfg)
	ensure.Nil(t, err)

	send := []byte("YO")
	expect := []byte("OK")

	servAddr := "localhost:" + port
	ensure.Nil(t, err)

	resp, err := testutil.SendTCP(servAddr, send)
	ensure.Nil(t, err)
	ensure.DeepEqual(t, resp, expect)

	// re-activate, make sure it's a no-op
	err = Activate(cfg, cfg)
	ensure.Nil(t, err)

	resp, err = testutil.SendTCP(servAddr, send)
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
