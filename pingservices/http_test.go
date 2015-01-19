package pingservices

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestSuccessfulPing(t *testing.T) {

	assert := assert.New(t)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, client")
	}))
	defer ts.Close()

	http := new(HTTP)
	result := http.Ping(ts.URL, time.Second*10)

	assert.Equal(result.StatusCode, 200, "should return 200")
}

func TestRedirectStatusCodePing(t *testing.T) {

	assert := assert.New(t)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/someotherurl", http.StatusFound)
	}))
	defer ts.Close()

	http := new(HTTP)
	result := http.Ping(ts.URL, time.Second*10)

	assert.Equal(result.StatusCode, 302, "should return 302")
}

// TODO: a test for timeout
