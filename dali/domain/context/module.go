package domain_context

import (
	"github.com/net12labs/cirm/dali/context/cmd"
	webserver "github.com/net12labs/cirm/mali/web-server"
)

type DomainUnit struct {
	Path        string
	Mode        RunMode
	OnExit      func()
	ExitMessage string
	WebServer   *webserver.WebServer
	Execute     func(cmd *cmd.Cmd)

	// Other fields here
}

type Domain struct {
	*DomainUnit

	// Other fields here
}

type SubDomain struct {
	*DomainUnit

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

func (a *DomainUnit) OnExecute(cmd *cmd.Cmd) {
	a.Execute(cmd)
}
