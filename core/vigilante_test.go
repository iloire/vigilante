package core

import (
	"testing"
)

func TestVigilanteAddService(t *testing.T) {
	v := Vigilante{}
	v.AddService(Service{"service 1", "http://google.com", 5000, new(MockPingService)})
	v.AddService(Service{"service 2", "http://yahoo.com", 4000, new(MockPingService)})

	if len(v.GetServices()) != 2 {
		t.Error("We are supossed to have one element")
	}

	if v.GetServices()[0].Name != "service 1" {
		t.Error("Wrong service 1")
	}

	if v.GetServices()[1].Name != "service 2" {
		t.Error("Wrong service 1")
	}
}
