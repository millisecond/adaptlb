package testutil

import (
	"fmt"
	"github.com/facebookgo/ensure"
	"io/ioutil"
	"net/http"
	"strconv"
	"testing"
)

func TestHTTPServer(port int) *http.Server {
	handler := func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprint(w, "OK")
	}
	srv := &http.Server{
		Addr:    "localhost:" + strconv.Itoa(port),
		Handler: http.HandlerFunc(handler),
	}
	go func() {
		srv.ListenAndServe()
	}()
	return srv
}

func MustBody(t *testing.T, resp *http.Response) string {
	b, err := ioutil.ReadAll(resp.Body)
	ensure.Nil(t, err)
	return string(b)
}
