package rules

import (
	"strings"
)

type Contains struct {
	Content string
}

func (c Contains) Match(Response string, StatusCode int) MatchResult {
	if strings.Contains(Response, c.Content) {
		return MatchResult{true, ""}
	} else {
		return MatchResult{false, "Response doesn't contain " + c.Content}
	}
}
