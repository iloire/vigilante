package pingservices

import (
	"fmt"
	"time"
)

type PingResult struct {
	Success bool
	Elapsed time.Duration
	Errors  []string
}

func (p *PingResult) Log() {
	fmt.Printf("success: %t\n", p.Success)
	fmt.Printf("elapsed: %v\n", p.Elapsed)
}
