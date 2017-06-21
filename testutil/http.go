package testutil

import (
	"github.com/facebookgo/ensure"
	"io/ioutil"
	"net/http"
	"testing"
)

func MustBody(t *testing.T, resp *http.Response) string {
	b, err := ioutil.ReadAll(resp.Body)
	ensure.Nil(t, err)
	return string(b)
}
