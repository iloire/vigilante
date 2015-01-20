package pingservices

import (
	"time"
	"vigilante/rules"
)

type PingService interface {
	Ping(url string, timeout time.Duration, rules []rules.Rule) PingResult
}
