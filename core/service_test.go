package core

import (
	"github.com/iloire/vigilante/pingservices"
	"github.com/iloire/vigilante/rules"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

// Mock Clock for testing
type MockClock struct {
	time time.Time
}

func (m *MockClock) SetTime(time time.Time) {
	m.time = time
}
func (m *MockClock) Forward(d time.Duration) {
	m.time = m.time.Add(d)
}
func (m MockClock) Now() time.Time                  { return m.time }
func (m MockClock) Since(t time.Time) time.Duration { return m.time.Sub(t) }

// Mock Ping Service
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
	service, _, _ := getServiceForTesting(10 * time.Millisecond)

	c := make(chan pingservices.PingResult)
	go service.Start(c)

	<-c

	assert.True(service.IsEnabled())
	service.Stop()
	assert.False(service.IsEnabled())
}

func TestServiceInterval(t *testing.T) {

	assert := assert.New(t)
	service, _, mockPingService := getServiceForTesting(10 * time.Millisecond)
	mockPingService.SetNextResult(pingservices.PingResult{true, 2, nil})

	go service.Start(nil)

	time.Sleep(20 * time.Millisecond)

	service.Stop()

	assert.Equal(2, service.GetTotalCount())
	assert.Equal(2, service.GetSuccessCount())
}

func TestServiceRecoveryInterval(t *testing.T) {

	assert := assert.New(t)
	service, _, mockPingService := getServiceForTesting(10 * time.Millisecond)
	service.RecoveryInterval = 2
	mockPingService.SetNextResult(pingservices.PingResult{false, 5, []string{"Error"}})

	go service.Start(nil)

	time.Sleep(20 * time.Millisecond)

	service.Stop()

	assert.Equal(3, service.GetTotalCount())
	assert.Equal(3, service.GetErrorCount())
}

func TestServiceUndefinedRecoveryIntervalDefaultsToInterval(t *testing.T) {

	assert := assert.New(t)

	service, _, mockPingService := getServiceForTesting(10 * time.Millisecond)
	mockPingService.SetNextResult(pingservices.PingResult{false, 2, []string{"Error"}})

	go service.Start(nil)

	time.Sleep(20 * time.Millisecond)

	service.Stop()

	assert.Equal(2, service.GetTotalCount())
	assert.Equal(2, service.GetErrorCount())
}

func TestRunningTotalTime(t *testing.T) {

	assert := assert.New(t)

	service, mockClock, mockPingService := getServiceForTesting(10 * time.Millisecond)
	mockPingService.SetNextResult(pingservices.PingResult{true, 5, nil})

	c := make(chan pingservices.PingResult)
	go service.Start(c)

	time.Sleep(time.Millisecond * 5)

	<-c
	time.Sleep(time.Millisecond * 15)
	mockClock.Forward(1 * time.Second)
	assert.Equal(0, service.GetRunningTotal()) // the first ping doesn't have previous reference

	<-c
	assert.Equal(1*time.Second, service.GetRunningTotal())

	service.Stop()
}

func TestUpTime(t *testing.T) {

	assert := assert.New(t)

	service, mockClock, mockPingService := getServiceForTesting(10 * time.Millisecond)
	mockPingService.SetNextResult(pingservices.PingResult{false, 5, []string{"Error"}})

	c := make(chan pingservices.PingResult)
	go service.Start(c)

	time.Sleep(time.Millisecond * 5) // so the first ping can be issued before setting the next ping result

	// prepare second ping
	mockPingService.SetNextResult(pingservices.PingResult{true, 10, nil})

	<-c                                // result from first ping
	mockClock.Forward(1 * time.Second) // mock time for second ping

	// assert
	assert.Equal(0, service.GetUpTime(), "0% uptime expected")
	assert.Equal(0, service.GetRunningTotal()) // first ping
	assert.Equal(1, service.GetTotalCount())

	mockPingService.SetNextResult(pingservices.PingResult{false, 5, []string{"Error"}})

	<-c // result from second ping
	assert.Equal(100, service.GetUpTime(), "100% uptime expected")
	assert.Equal(2, service.GetTotalCount())
	assert.Equal(1*time.Second, service.GetRunningTotal())

	mockClock.Forward(1 * time.Second) // mock clock for third ping

	<-c // result from third ping
	assert.Equal(50, service.GetUpTime(), "50% uptime expected")
	assert.Equal(3, service.GetTotalCount())

	mockClock.Forward(1 * time.Second) // mock clock for 4th ping

	<-c // result from 4th ping
	assert.Equal(33, service.GetUpTime(), "33% uptime expected")
	assert.Equal(4, service.GetTotalCount())

	service.Stop()
}

// helpers
func getServiceForTesting(interval time.Duration) (*Service, *MockClock, *MockPingService) {
	mockedClock := new(MockClock)
	mockedClock.SetTime(time.Date(2000, time.November, 1, 1, 0, 0, 0, time.UTC))

	mockedPingService := new(MockPingService)

	service := &Service{
		Name:        "service google",
		Url:         "http://google.com",
		Interval:    interval,
		Timeout:     10 * time.Second,
		PingService: mockedPingService,
		Clock:       mockedClock}

	return service, mockedClock, mockedPingService
}
