package main

import (
	"github.com/guionardo/gs-service-scheduler/logging"
	"github.com/guionardo/gs-service-scheduler/service"
	"github.com/guionardo/gs-service-scheduler/setup"
)

func main() {
	args, err := setup.GetArgs()
	if err != nil {
		logging.ErrorF("Argument error - %v", err)
		return
	}
	svc := service.CreateService(service.Runner, args)
	service.RunService(svc)
}
