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
		return PingResult{500, elapsed, false}
	}

	defer resp.Body.Close()
	ioutil.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return PingResult{resp.StatusCode, elapsed, false}
	}
	return PingResult{200, elapsed, true}
}
