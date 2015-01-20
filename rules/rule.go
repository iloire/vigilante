package rules

type MatchResult struct {
	Success bool
	Error   string
}

type Rule interface {
	Match(Response string, StatusCode int) MatchResult
}
