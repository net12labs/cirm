package service

import webserver "github.com/net12labs/cirm/dali/web-server"

type ServiceUnit struct {
	Mode        RunMode
	OnExit      func()
	WebServer   *webserver.WebServer
	ExitMessage string
	// Other fields here
}

type Service struct {
	ServiceUnit
	// Other fields here
}

type SubService struct {
	ServiceUnit
	// Other fields here
}

func NewService() *Service {
	svc := &Service{}
	svc.ServiceUnit = ServiceUnit{
		Mode: RunMode{
			items: make(map[string]string),
		},
	}
	return svc
}

func NewSubService() *SubService {
	svc := &SubService{}
	svc.ServiceUnit = ServiceUnit{
		Mode: RunMode{
			items: make(map[string]string),
		},
	}
	return svc
}

var NewWebServer = webserver.NewWebServer
