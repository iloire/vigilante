package core

import (
	"fmt"
	"testing"
	"time"
	"vigilante/pingservices"
)

type MockPingService struct{}

func (m *MockPingService) Ping(url string) pingservices.PingResult {
	return pingservices.PingResult{200, 100, true}
}

//TODO: investigate how to write shorter assertions
//TODO: write more unit (less general) tests
func TestService(t *testing.T) {
	service := Service{
		"service 1",
		"http://google.com",
		100,
		new(MockPingService),
	}

	if service.Url != "http://google.com" {
		t.Error("Invalid url")
	}

	go func(service Service) {
		fmt.Printf(service.Name)
		service.Start()
	}(service)

	time.Sleep(time.Second * 1) // we need to somehow fake time

	if service.IsEnabled() != true {
		t.Error("It should be enabled")
	}

	service.Stop()

	if service.IsEnabled() != false {
		t.Error("It should be disabled")
	}

	if service.GetTotalCount() != 10 {
		t.Error("Incorrect total count")
	}

	if service.GetSuccessCount() != 10 {
		t.Error("Incorrect success count")
	}

	if service.GetErrorCount() != 0 {
		t.Error("Incorrect error count")
	}

	if service.GetAVGLatency() != 100 {
		t.Error("Incorrect latency")
	}
}
