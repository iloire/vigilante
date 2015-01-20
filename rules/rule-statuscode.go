package rules

import ()

type StatusCode struct {
	StatusCode int
}

func (c StatusCode) Match(Response string, StatusCode int) MatchResult {
	if c.StatusCode != StatusCode {
		return MatchResult{false, "Invalid status code"}
	} else {
		return MatchResult{true, ""}
	}
}
