package core

import (
	"github.com/iloire/vigilante/pingservices"
	"github.com/iloire/vigilante/rules"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type MockPingService struct{}

func (m *MockPingService) Ping(url string, timeout time.Duration, rules []rules.Rule) pingservices.PingResult {
	return pingservices.PingResult{true, 100, nil}
}

func TestService(t *testing.T) {

	assert := assert.New(t)

	service := Service{
		Name:        "service google",
		Url:         "http://google.com",
		Interval:    100,
		Timeout:     10000,
		PingService: new(MockPingService)}

	assert.Equal(service.Url, "http://google.com", "valid url")

	go func(service *Service) {
		service.Start()
	}(&service)

	time.Sleep(time.Second * 1) // we need to somehow fake time

	assert.Equal(service.IsEnabled(), true, "service is enable")

	service.Stop()

	assert.Equal(service.IsEnabled(), false, "service is disable")
	assert.Equal(service.GetTotalCount(), 10, "total count")
	assert.Equal(service.GetSuccessCount(), 10, "success count")
	assert.Equal(service.GetErrorCount(), 0, "error count")
	assert.Equal(service.GetAVGLatency(), 100, "avg latency")
}
