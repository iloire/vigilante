package pingservices

import (
	"github.com/iloire/vigilante/rules"
	"time"
)

type PingService interface {
	Ping(url string, timeout time.Duration, rules []rules.Rule) PingResult
}
