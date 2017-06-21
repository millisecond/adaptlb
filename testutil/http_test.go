package testutil

import (
	"testing"
	"net/http"
	"github.com/facebookgo/ensure"
)

func TestStoppableHTTPListener(t *testing.T) {
	srv := TestHTTPServer(8001)
	url := "http://localhost:8001"

	resp, err := http.Get(url)
	ensure.Nil(t, err)
	ensure.DeepEqual(t, resp.StatusCode, 200)

	err = srv.Close()
	ensure.Nil(t, err)

	resp, err = http.Get(url)
	ensure.NotNil(t, err)
}
