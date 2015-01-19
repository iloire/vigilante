package pingservices

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type HTTP struct{}

func (h *HTTP) Ping(url string, timeout time.Duration) PingResult {

	// TODO implement timeout (not straightforward)

	start := time.Now()
	resp, err := http.Get(url)
	elapsed := time.Since(start)

	if err != nil {
		fmt.Println(err)
		return PingResult{resp.StatusCode, elapsed}
	}

	// TODO: eventually we need to be able to just ping to get the headers only
	// so we don't download the entire page.
	// download the entire page will be necessary if we want to assert for a certain content though
	defer resp.Body.Close()
	ioutil.ReadAll(resp.Body)

	return PingResult{resp.StatusCode, elapsed}
}
