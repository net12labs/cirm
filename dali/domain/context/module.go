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
}

type Domain struct {
	*DomainUnit
	ApiRequest  func(req *webserver.Request)
	PageRequest func(req *webserver.Request)
	WebRequest  func(req *webserver.Request)
}

type SubDomain struct {
	*DomainUnit
	ApiRequest  func(req *webserver.Request)
	PageRequest func(req *webserver.Request)
	WebRequest  func(req *webserver.Request)
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

func (a *SubDomain) OnApiRequest(req *webserver.Request) {
	a.ApiRequest(req)
}

func (a *Domain) OnApiRequest(req *webserver.Request) {
	a.ApiRequest(req)
}

func (a *SubDomain) OnWebRequest(req *webserver.Request) {
	a.WebRequest(req)
}

func (a *Domain) OnWebRequest(req *webserver.Request) {
	a.WebRequest(req)
}
