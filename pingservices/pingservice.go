package pingservices

import (
	"time"
)

type PingService interface {
	Ping(url string, timeout time.Duration) PingResult
}
