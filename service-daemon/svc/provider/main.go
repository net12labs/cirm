package provider

import (
	"github.com/net12labs/cirm/dali/context/service"
	webserver "github.com/net12labs/cirm/dali/web-server"

	webclient "github.com/net12labs/cirm/client-web/provider"
)

// Possible runmodes are; web, cli

type Unit struct {
	*service.Service
	WebServer   *webserver.WebServer
	Webclient   *webclient.ProviderClient
	Agent       *SvcAgent
	ExitMessage string
	WebApi      *WebApi
}

func NewUnit() *Unit {
	svc := &Unit{}
	svc.Service = service.NewService()
	svc.Webclient = webclient.NewProviderClient()
	svc.WebApi = NewWebApi()
	svc.WebApi.svc = svc
	svc.Agent = &SvcAgent{Svc: svc}
	return svc
}

func (r *Unit) Init() error {
	r.Webclient.Server = r.WebServer
	r.Webclient.Init()
	r.WebApi.Server = r.WebServer
	r.WebApi.Init()

	return nil
}

func (r *Unit) Run() int {

	// Initialize othe components here

	// Initialize other components here
	// Start the application
	return 0
}
