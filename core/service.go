package core

import (
	"fmt"
	"github.com/iloire/vigilante/pingservices"
	"github.com/iloire/vigilante/rules"
	"time"
)

type Service struct {
	Name             string
	Url              string
	Interval         time.Duration
	RecoveryInterval time.Duration
	PingService      pingservices.PingService
	Timeout          time.Duration
	Rules            []rules.Rule
	Clock            Clock

	runningTotal                 time.Duration
	errorTotal                   time.Duration
	lastResult                   *LastResult
	enabled                      bool
	totalcounter                 int
	successcounter, errorcounter int
	avgLatency                   time.Duration
}

type RealClock struct{}

func (m RealClock) Now() time.Time                  { return time.Now() }
func (m RealClock) Since(t time.Time) time.Duration { return time.Since(t) }

type LastResult struct {
	TimeStamp  time.Time
	PingResult pingservices.PingResult
}

// Starts the service.
// It will be executing the "PingService" every "Interval" until it gets stopped.
func (s *Service) Start(c chan pingservices.PingResult) {

	if s.Clock == nil {
		s.Clock = new(RealClock)
	}

	s.enabled = true

	for s.enabled {
		result := s.PingService.Ping(s.Url, s.Timeout, s.Rules)
		s.avgLatency = time.Duration((int(s.avgLatency)*s.totalcounter + int(result.Elapsed)) / (s.totalcounter + 1))

		s.totalcounter++

		var nextInterval = s.Interval

		if s.lastResult != nil {
			s.runningTotal += s.Clock.Since(s.lastResult.TimeStamp)
			if !s.lastResult.PingResult.Success {
				s.errorTotal += s.Clock.Since(s.lastResult.TimeStamp)
			}
		}

		if result.Success {
			s.successcounter++
		} else {
			s.errorcounter++
			if s.RecoveryInterval != 0 {
				nextInterval = s.RecoveryInterval
			}
		}

		s.lastResult = &LastResult{s.Clock.Now(), result}

		if c != nil {
			c <- result
		}

		time.Sleep(nextInterval)
	}
}

// Stops the service
func (s *Service) Stop() {
	s.enabled = false
}

func (s *Service) Log() {
	fmt.Printf("--\nPing from %s: success: %d, error: %d\n", s.Name, s.successcounter, s.errorcounter)
}

func (s *Service) GetLastResult() *LastResult {
	return s.lastResult
}

func (s *Service) GetUpTime() int {
	if s.runningTotal == 0 {
		return 0
	}
	return (int)((s.runningTotal*time.Nanosecond - s.errorTotal*time.Nanosecond) * 100 / s.runningTotal * time.Nanosecond)
}

func (s *Service) GetRunningTotal() time.Duration {
	return s.runningTotal
}

func (s *Service) GetErrorTotal() time.Duration {
	return s.errorTotal
}

// Number of ping executions
func (s *Service) GetTotalCount() int {
	return s.totalcounter
}

// Number of times the service is alive
func (s *Service) GetSuccessCount() int {
	return s.successcounter
}

// Number of times the server was down or not reacheable
func (s *Service) GetErrorCount() int {
	return s.errorcounter
}

func (s *Service) GetAVGLatency() time.Duration {
	return s.avgLatency
}

func (s *Service) IsEnabled() bool {
	return s.enabled
}
