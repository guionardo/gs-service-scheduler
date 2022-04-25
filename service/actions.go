package service

import (
	"fmt"
	"log"

	"github.com/guionardo/gs-service-scheduler/setup"
	"github.com/kardianos/service"
)

type Service struct {
	RunnerMethod func(service.Service, *setup.Args) error
	Args         *setup.Args
}

func (p Service) Start(s service.Service) error {
	fmt.Println(s.String() + " started")
	args, err := setup.GetArgs()
	if err != nil {
		return err
	}
	go p.Run(s, args)

	return nil
}

func (p Service) Stop(s service.Service) error {
	fmt.Println(s.String() + " stopped")
	return nil
}

func (p Service) Run(s service.Service, args *setup.Args) error {
	if p.RunnerMethod == nil {
		p.RunnerMethod = func(svc service.Service, args *setup.Args) error {
			log.Printf("DUMMY RUNNER FOR SERVICE %s - %v", svc, args)
			return nil
		}
	}
	return p.RunnerMethod(s, args)
}
