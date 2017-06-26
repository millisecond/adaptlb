package testutil

import (
	"github.com/facebookgo/ensure"
	"net/http"
	"strconv"
	"testing"
)

func TestStoppableHTTPListener(t *testing.T) {
	t.Parallel()

	port := UniquePort()
	srv := TestHTTPServer(port)
	url := "http://localhost:" + strconv.Itoa(port)

	resp, err := http.Get(url)
	ensure.Nil(t, err)
	ensure.DeepEqual(t, resp.StatusCode, 200)

	err = srv.Close()
	ensure.Nil(t, err)

	resp, err = http.Get(url)
	ensure.NotNil(t, err)
}
