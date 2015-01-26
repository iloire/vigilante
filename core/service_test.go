package core

import (
	"github.com/iloire/vigilante/pingservices"
	"github.com/iloire/vigilante/rules"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type MockPingService struct {
	nextPingResult pingservices.PingResult
}

func (m *MockPingService) Ping(url string, timeout time.Duration, rules []rules.Rule) pingservices.PingResult {
	return m.nextPingResult
}

func (m *MockPingService) SetNextResult(result pingservices.PingResult) {
	m.nextPingResult = result
}

func TestServiceIsEnabled(t *testing.T) {

	assert := assert.New(t)

	mockPingService := new(MockPingService)
	mockPingService.SetNextResult(pingservices.PingResult{true, 100, nil})

	service := Service{
		Name:        "service google",
		Url:         "http://google.com",
		Interval:    1,
		Timeout:     10000,
		PingService: mockPingService}

	c := make(chan pingservices.PingResult)
	go func(service *Service) {
		service.Start(c)
	}(&service)

	<-c
	assert.True(service.IsEnabled(), "service is enable")
	service.Stop()
	assert.False(service.IsEnabled(), "service is disable")
}

func TestServiceInterval(t *testing.T) {

	assert := assert.New(t)

	mockPingService := new(MockPingService)
	mockPingService.SetNextResult(pingservices.PingResult{true, 100, nil})

	service := Service{
		Name:        "service google",
		Url:         "http://google.com",
		Interval:    1,
		Timeout:     10000,
		PingService: mockPingService}

	c := make(chan pingservices.PingResult)
	go func(service *Service) {
		service.Start(c)
	}(&service)

	<-c
	<-c

	service.Stop()

	assert.Equal(2, service.GetTotalCount(), "total count")
	assert.Equal(2, service.GetSuccessCount(), "success count")
	assert.Equal(0, service.GetErrorCount(), "error count")
	assert.Equal(100, service.GetAVGLatency(), "avg latency")
}

func TestServiceRecoveryInterval(t *testing.T) {

	assert := assert.New(t)

	mockPingService := new(MockPingService)
	mockPingService.SetNextResult(pingservices.PingResult{false, 200, []string{"Error"}})

	service := Service{
		Name:             "service google",
		Url:              "http://google.com",
		Interval:         10,
		RecoveryInterval: 5,
		Timeout:          10000,
		PingService:      mockPingService}

	go func(service *Service) {
		service.Start(nil)
	}(&service)

	time.Sleep(time.Millisecond * 20) // we need to somehow fake time

	service.Stop()

	assert.Equal(4, service.GetTotalCount(), "total count")
	assert.Equal(0, service.GetSuccessCount(), "success count")
	assert.Equal(4, service.GetErrorCount(), "error count")
}

func TestServiceRecoveryIntervalDefaultsToInterval(t *testing.T) {

	assert := assert.New(t)

	mockPingService := new(MockPingService)
	mockPingService.SetNextResult(pingservices.PingResult{false, 200, []string{"Error"}})

	service := Service{
		Name:        "service google",
		Url:         "http://google.com",
		Interval:    10,
		Timeout:     10000,
		PingService: mockPingService}

	go func(service *Service) {
		service.Start(nil)
	}(&service)

	time.Sleep(time.Millisecond * 20) // we need to somehow fake time

	service.Stop()

	assert.Equal(2, service.GetTotalCount(), "total count")
	assert.Equal(0, service.GetSuccessCount(), "success count")
	assert.Equal(2, service.GetErrorCount(), "error count")
}
