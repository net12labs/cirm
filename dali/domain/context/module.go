package domain_context

import (
	webserver "github.com/net12labs/cirm/mali/web-server"
)

type DomainUnit struct {
	Path        string
	Mode        RunMode
	OnExit      func()
	ExitMessage string
	WebServer   *webserver.WebServer
	ApiRequest  func(req *webserver.Request)

	// Other fields here
}

type Domain struct {
	*DomainUnit
	PageRequest func(req *webserver.Request)
	// Other fields here
}

type SubDomain struct {
	*DomainUnit
	PageRequest func(req *webserver.Request)
	// Other fields here
}

func NewDomain() *Domain {
	d := &Domain{}
	d.DomainUnit = &DomainUnit{
		Mode: RunMode{
			items: make(map[string]string),
		},
	}
	return d
}

func NewSubDomain() *SubDomain {
	svc := &SubDomain{}
	svc.DomainUnit = &DomainUnit{
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

func (su *DomainUnit) SetPath(key string) {
	su.Path = key
}

func (a *DomainUnit) OnApiRequest(req *webserver.Request) {
	a.ApiRequest(req)
}
