package service

import (
	"github.com/net12labs/cirm/dali/context/cmd"
	webserver "github.com/net12labs/cirm/mali/web-server"
)

type ServiceUnit struct {
	Path        string
	Mode        RunMode
	OnExit      func()
	ExitMessage string
	WebServer   *webserver.WebServer
	Execute     func(cmd *cmd.Cmd)
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

func NewServer() *webserver.WebServer {
	return webserver.NewWebServer()
}

type WebServer = webserver.WebServer

func (su *ServiceUnit) SetPath(key string) {
	su.Path = key
}
