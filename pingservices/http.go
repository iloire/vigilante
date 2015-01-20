package pingservices

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
	"vigilante/rules"
)

type HTTP struct{}

func (h *HTTP) Ping(url string, timeout time.Duration, rules []rules.Rule) PingResult {

	// TODO implement timeout (not straightforward)

	start := time.Now()
	resp, err := http.Get(url)
	elapsed := time.Since(start)

	if err != nil {
		fmt.Println(err)
		return PingResult{true, elapsed, []string{err.Error()}}
	}

	// TODO: eventually we need to be able to just ping to get the headers only
	// so we don't download the entire page.
	// download the entire page will be necessary if we want to assert for a certain content though
	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return PingResult{false, elapsed, []string{err.Error()}}
	}

	for _, rule := range rules {
		result := rule.Match(string(content), resp.StatusCode)
		if !result.Success {
			return PingResult{false, elapsed, []string{result.Error}}
		}
	}

	return PingResult{true, elapsed, nil}
}
