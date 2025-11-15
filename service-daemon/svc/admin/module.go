package admin

import (
	"github.com/net12labs/cirm/dali/context/service"
	webserver "github.com/net12labs/cirm/dali/web-server"

	webclient "github.com/net12labs/cirm/client-web/admin"
)

// Possible runmodes are; web, cli

type Unit struct {
	*service.Service
	WebServer   *webserver.WebServer
	WebClient   *webclient.AdminClient
	Agent       *SvcAgent
	ExitMessage string
	WebApi      *WebApi
}

func NewUnit() *Unit {
	svc := &Unit{}
	svc.Service = service.NewService()
	svc.WebClient = webclient.NewAdminClient()
	svc.WebApi = NewWebApi()
	svc.WebApi.svc = svc
	svc.Agent = &SvcAgent{Svc: svc}

	return svc
}

func (r *Unit) Init() error {
	r.WebClient.Server = r.WebServer
	r.WebClient.Init()
	r.WebApi.Server = r.WebServer
	r.WebApi.Init()
	return nil
}

func (r *Unit) Run() int {

	// Initialize other components here
	// Start the application
	return 0
}
