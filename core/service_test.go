package core

import (
	"github.com/iloire/vigilante/pingservices"
	"github.com/iloire/vigilante/rules"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type MockClock struct {
	time time.Time
}

func (m *MockClock) SetTime(time time.Time) {
	m.time = time
}
func (m MockClock) Now() time.Time                  { return m.time }
func (m MockClock) Since(t time.Time) time.Duration { return m.time.Sub(t) }

type MockPingService struct {
	nextPingResult pingservices.PingResult
}

func (m *MockPingService) Ping(url string, timeout time.Duration, rules []rules.Rule) pingservices.PingResult {
	time.Sleep(m.nextPingResult.Elapsed * time.Millisecond)
	return m.nextPingResult
}

func (m *MockPingService) SetNextResult(result pingservices.PingResult) {
	m.nextPingResult = result
}

func TestServiceIsEnabled(t *testing.T) {

	assert := assert.New(t)

	mockPingService := new(MockPingService)
	mockPingService.SetNextResult(pingservices.PingResult{true, 1, nil})

	mockClock := new(MockClock)
	now := time.Date(2000, time.November, 1, 1, 0, 0, 0, time.UTC)
	mockClock.SetTime(now)

	service := Service{
		Name:        "service google",
		Url:         "http://google.com",
		Interval:    1,
		Timeout:     10000,
		PingService: mockPingService,
		Clock:       mockClock}

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

	mockClock := new(MockClock)
	now := time.Date(2000, time.November, 1, 1, 0, 0, 0, time.UTC)
	mockClock.SetTime(now)

	service := Service{
		Name:        "service google",
		Url:         "http://google.com",
		Interval:    1,
		Timeout:     10000,
		PingService: mockPingService,
		Clock:       mockClock}

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
	mockPingService.SetNextResult(pingservices.PingResult{false, 5, []string{"Error"}})

	mockClock := new(MockClock)
	now := time.Date(2000, time.November, 1, 1, 0, 0, 0, time.UTC)
	mockClock.SetTime(now)

	service := Service{
		Name:             "service google",
		Url:              "http://google.com",
		Interval:         10,
		RecoveryInterval: 5,
		Timeout:          10000,
		PingService:      mockPingService,
		Clock:            mockClock}

	go func(service *Service) {
		service.Start(nil)
	}(&service)

	time.Sleep(time.Millisecond * 50) // we need to somehow fake time

	service.Stop()

	assert.Equal(4, service.GetTotalCount(), "total count")
	assert.Equal(0, service.GetSuccessCount(), "success count")
	assert.Equal(4, service.GetErrorCount(), "error count")
}

func TestServiceRecoveryIntervalDefaultsToInterval(t *testing.T) {

	assert := assert.New(t)

	mockPingService := new(MockPingService)
	mockPingService.SetNextResult(pingservices.PingResult{false, 5, []string{"Error"}})

	mockClock := new(MockClock)
	now := time.Date(2000, time.November, 1, 1, 0, 0, 0, time.UTC)
	mockClock.SetTime(now)

	service := Service{
		Name:        "service google",
		Url:         "http://google.com",
		Interval:    10,
		Timeout:     10000,
		PingService: mockPingService,
		Clock:       mockClock}

	go func(service *Service) {
		service.Start(nil)
	}(&service)

	time.Sleep(time.Millisecond * 30) // we need to somehow fake time

	service.Stop()

	assert.Equal(2, service.GetTotalCount(), "total count")
	assert.Equal(0, service.GetSuccessCount(), "success count")
	assert.Equal(2, service.GetErrorCount(), "error count")
}

func TestRunningTotal(t *testing.T) {

	assert := assert.New(t)

	mockPingService := new(MockPingService)
	mockPingService.SetNextResult(pingservices.PingResult{false, 5, []string{"Error"}})

	mockClock := new(MockClock)
	now := time.Date(2000, time.November, 1, 1, 0, 0, 0, time.UTC)
	mockClock.SetTime(now)

	service := Service{
		Name:        "service google",
		Url:         "http://google.com",
		Interval:    10,
		Timeout:     10000,
		PingService: mockPingService,
		Clock:       mockClock}

	c := make(chan pingservices.PingResult)
	go func(service *Service) {
		service.Start(c)
	}(&service)

	<-c
	assert.Equal(0, service.GetRunningTotal())

	now = now.Add(1 * time.Second)
	mockClock.SetTime(now) // mock time for second ping

	<-c

	assert.Equal(1*time.Second, service.GetRunningTotal())

	service.Stop()
}

func TestUpTime(t *testing.T) {

	assert := assert.New(t)

	mockPingService := new(MockPingService)
	mockPingService.SetNextResult(pingservices.PingResult{false, 5, []string{"Error"}})

	mockClock := new(MockClock)
	now := time.Date(2000, time.November, 1, 1, 0, 0, 0, time.UTC)
	mockClock.SetTime(now)

	service := Service{
		Name:        "service google",
		Url:         "http://google.com",
		Interval:    10,
		Timeout:     10000,
		PingService: mockPingService,
		Clock:       mockClock}

	c := make(chan pingservices.PingResult)
	go func(service *Service) {
		service.Start(c)
	}(&service)

	time.Sleep(time.Millisecond * 5) // so the first ping can be issued before setting the next ping result

	// prepare second ping
	mockPingService.SetNextResult(pingservices.PingResult{true, 10, nil})

	<-c // result from first ping
	now = now.Add(1 * time.Second)
	mockClock.SetTime(now) // mock time for second ping

	// assert
	assert.Equal(0, service.GetUpTime(), "0% uptime expected")
	assert.Equal(0, service.GetRunningTotal()) // first ping
	assert.Equal(1, service.GetTotalCount())

	mockPingService.SetNextResult(pingservices.PingResult{false, 5, []string{"Error"}})

	<-c // result from second ping
	assert.Equal(100, service.GetUpTime(), "100% uptime expected")
	assert.Equal(2, service.GetTotalCount())
	assert.Equal(1*time.Second, service.GetRunningTotal())

	now = now.Add(1 * time.Second) // mock clock for third ping
	mockClock.SetTime(now)

	<-c // result from third ping
	assert.Equal(50, service.GetUpTime(), "50% uptime expected")
	assert.Equal(3, service.GetTotalCount())

	now = now.Add(1 * time.Second) // mock clock for 4th ping
	mockClock.SetTime(now)

	<-c // result from 4th ping
	assert.Equal(33, service.GetUpTime(), "33% uptime expected")
	assert.Equal(4, service.GetTotalCount())

	service.Stop()
}
