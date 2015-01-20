package pingservices

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"vigilante/rules"
)

func TestHTTPPingSuccessfulPingNoRules(t *testing.T) {

	assert := assert.New(t)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, client")
	}))
	defer ts.Close()

	http := new(HTTP)
	result := http.Ping(ts.URL, time.Second*10, []rules.Rule{})

	assert.True(result.Success, "should be successful")
}

func TestHTTPPingStatusCodeRule(t *testing.T) {

	assert := assert.New(t)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/someotherurl", http.StatusNotFound)
	}))
	defer ts.Close()

	http := new(HTTP)
	result := http.Ping(ts.URL, time.Second*10, []rules.Rule{rules.StatusCode{StatusCode: 200}})

	assert.False(result.Success, "should error")
}

// TODO: a test for timeout
