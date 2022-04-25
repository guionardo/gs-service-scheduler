package service

import (
	"log"

	"github.com/guionardo/gs-service-scheduler/setup"
	"github.com/kardianos/service"
)

const serviceName = "GS Task Scheduler"
const serviceDescription = "Multiple tasks scheduler"

func CreateService(runnerMethod func(service.Service)error, args *setup.Args) service.Service {
	serviceConfig := &service.Config{
		Name:        serviceName,
		DisplayName: serviceName,
		Description: serviceDescription,
	}
	prg := &Service{}
	s, err := service.New(prg, serviceConfig)
	if err != nil {
		log.Fatalf("Cannot create the service: %v", err)
	}
	return s

}

func RunService(svc service.Service) {
	if err := svc.Run(); err != nil {
		log.Fatalf("Cannot start the service: %v", err)
	}
}
