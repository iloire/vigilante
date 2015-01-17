package pingservices

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	// "time"
)

func TestSuccessfulPing(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//		time.Sleep(time.Millisecond * 10)
		fmt.Fprintln(w, "Hello, client")
	}))
	defer ts.Close()

	http := new(HTTP)
	result := http.Ping(ts.URL)
	if result.Status != 200 {
		t.Errorf("Invalid response")
	}

	// if result.Elapsed != (time.Millisecond * 10) {
	// 	t.Errorf("Invalid duration %d", result.Elapsed)
	// }
}

// TODO: a test for timeout

// TODO: a test for invalid response
