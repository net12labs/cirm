package service

import (
	webserver "github.com/net12labs/cirm/dali/web-server"
)

type ServiceUnit struct {
	Mode        RunMode
	OnExit      func()
	ExitMessage string
	WebServer   *webserver.WebServer
	// Other fields here
}

type Service struct {
	*ServiceUnit

	// Other fields here
}

type SubService struct {
	*ServiceUnit

	// Other fields here
}

func NewService() *Service {
	svc := &Service{}
	svc.ServiceUnit = &ServiceUnit{
		Mode: RunMode{
			items: make(map[string]string),
		},
	}
	return svc
}

func NewSubService() *SubService {
	svc := &SubService{}
	svc.ServiceUnit = &ServiceUnit{
		Mode: RunMode{
			items: make(map[string]string),
		},
	}
	return svc
}

var NewServer = webserver.NewWebServer
