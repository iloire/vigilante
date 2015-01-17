package main

import (
	"sync"
	"vigilante/core"
	"vigilante/pingservices"
)

var wg sync.WaitGroup

func main() {

	v := core.Vigilante{}
	wg.Add(1)

	// TODO: this need to be read form a persistent datasource
	v.AddService(core.Service{
		Name:        "service google",
		Url:         "http://google.com",
		Interval:    5000,
		Timeout:     10000,
		PingService: new(pingservices.HTTP)})

	v.AddService(core.Service{
		Name:        "service yahoo",
		Url:         "http://yahoo.com",
		Interval:    5000,
		Timeout:     10000,
		PingService: new(pingservices.HTTP)})

	v.AddService(core.Service{
		Name:        "service localhost",
		Url:         "http://iloire.dyn.syd.atlassian.com:8080/confluence",
		Interval:    5000,
		Timeout:     5000,
		PingService: new(pingservices.HTTP)})

	v.AddService(core.Service{
		Name:        "service INVALID",
		Url:         "http://invalid",
		Interval:    5000,
		Timeout:     1000,
		PingService: new(pingservices.HTTP)})

	v.Start()
	wg.Wait()
}
