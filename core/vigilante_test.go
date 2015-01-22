package core

import (
	"github.com/iloire/vigilante/rules"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestVigilanteAddService(t *testing.T) {

	assert := assert.New(t)

	v := Vigilante{}

	v.AddService(Service{
		Name:        "service google",
		Url:         "http://google.com",
		Interval:    5000,
		Timeout:     10000,
		PingService: new(MockPingService),
		Rules:       []rules.Rule{&rules.Contains{Content: "google"}}})

	v.AddService(Service{
		Name:        "service yahoo",
		Url:         "http://yahoo.com",
		Interval:    5000,
		Timeout:     10000,
		PingService: new(MockPingService),
		Rules:       []rules.Rule{&rules.Contains{Content: "yahoo"}}})

	assert.Equal(len(v.GetServices()), 2, "number of services added is correct")
	assert.Equal(v.GetServices()[0].Name, "service google", "correct service name")
	assert.Equal(v.GetServices()[1].Name, "service yahoo", "correct service name")

}
