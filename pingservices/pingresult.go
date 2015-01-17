package pingservices

import (
	"fmt"
	"time"
)

type PingResult struct {
	Status  int
	Elapsed time.Duration
	Success bool
}

func (p *PingResult) Log() {
	fmt.Printf("status: %d\n", p.Status)
	fmt.Printf("elapsed: %v\n", p.Elapsed)
}
