package core

import (
	"fmt"
)

type Vigilante struct{}

var services []Service

func (v *Vigilante) AddService(service Service) {
	services = append(services, service)
}

func (v *Vigilante) GetServices() []Service {
	return services
}

func (v *Vigilante) Start() {
	for _, service := range services {
		go func(service Service) {
			fmt.Printf(service.Name)
			service.Start()
		}(service)
	}
}
