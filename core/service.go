package core

import (
	"fmt"
	"time"
	"vigilante/pingservices"
)

type Service struct {
	Name        string
	Url         string
	Interval    time.Duration
	PingService pingservices.PingService
	Timeout     time.Duration

	enabled                      bool
	totalcounter                 int
	successcounter, errorcounter int
	avgLatency                   time.Duration
}

// Starts the service.
// It will be executing the "PingService" every "Interval" until it gets stopped.
func (s *Service) Start() {
	fmt.Println("Starting service: " + s.Name + "...")
	s.enabled = true

	for {

		if s.enabled {

			result := s.PingService.Ping(s.Url, s.Timeout)

			s.avgLatency = time.Duration((int(s.avgLatency)*s.totalcounter + int(result.Elapsed)) / (s.totalcounter + 1))

			s.totalcounter++

			if result.Success {
				s.successcounter++
			} else {
				s.errorcounter++
			}

			// logging
			s.Log()
			result.Log()
			fmt.Printf("avgLatency: %v\n", s.avgLatency)

			time.Sleep(s.Interval * time.Millisecond)
		}

	}
}

// Stops the service
func (s *Service) Stop() {
	s.enabled = false
	fmt.Println("Stopping service: " + s.Name + "...")
}

func (s *Service) Log() {
	fmt.Printf("--\nPing from %s: success: %d, error: %d\n", s.Name, s.successcounter, s.errorcounter)
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
