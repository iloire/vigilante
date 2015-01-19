package pingservices

import (
	"fmt"
	"time"
)

type PingResult struct {
	StatusCode int
	Elapsed    time.Duration
}

func (p *PingResult) Log() {
	fmt.Printf("status: %d\n", p.StatusCode)
	fmt.Printf("elapsed: %v\n", p.Elapsed)
}
