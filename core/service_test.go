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
		Interval:    10,
		Timeout:     10000,
		PingService: mockPingService}

	go func(service *Service) {
		service.Start()
	}(&service)

	time.Sleep(time.Millisecond * 1) // TODO: find a better way

	assert.Equal(service.IsEnabled(), true, "service is enable")
	service.Stop()
	assert.Equal(service.IsEnabled(), false, "service is disable")
}

func TestServiceInterval(t *testing.T) {

	assert := assert.New(t)

	mockPingService := new(MockPingService)
	mockPingService.SetNextResult(pingservices.PingResult{true, 100, nil})

	service := Service{
		Name:        "service google",
		Url:         "http://google.com",
		Interval:    10,
		Timeout:     10000,
		PingService: mockPingService}

	go func(service *Service) {
		service.Start()
	}(&service)

	time.Sleep(time.Millisecond * 20) // we need to somehow fake time

	service.Stop()

	assert.Equal(service.GetTotalCount(), 2, "total count")
	assert.Equal(service.GetSuccessCount(), 2, "success count")
	assert.Equal(service.GetErrorCount(), 0, "error count")
	assert.Equal(service.GetAVGLatency(), 100, "avg latency")
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
		service.Start()
	}(&service)

	time.Sleep(time.Millisecond * 20) // we need to somehow fake time

	service.Stop()

	assert.Equal(service.GetTotalCount(), 4, "total count")
	assert.Equal(service.GetSuccessCount(), 0, "success count")
	assert.Equal(service.GetErrorCount(), 4, "error count")
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
		service.Start()
	}(&service)

	time.Sleep(time.Millisecond * 20) // we need to somehow fake time

	service.Stop()

	assert.Equal(service.GetTotalCount(), 2, "total count")
	assert.Equal(service.GetSuccessCount(), 0, "success count")
	assert.Equal(service.GetErrorCount(), 2, "error count")
}
