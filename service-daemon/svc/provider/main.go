package provider

import (
	"fmt"

	webserver "github.com/net12labs/cirm/dali/web-server"
	"github.com/net12labs/cirm/dali/work/service"

	webclient "github.com/net12labs/cirm/svc-client-web"
)

// Possible runmodes are; web, cli

type Unit struct {
	*service.ServiceUnit
	WebServer *webserver.WebServer

	// Other root fields here
	Webclient   *webclient.WebClient
	Agent       *SvcAgent
	ExitMessage string
	WebApi      *WebApi
}

func NewUnit() *Unit {
	svc := &Unit{}
	svc.ServiceUnit = service.NewServiceUnit()
	svc.Webclient = webclient.NewWebClient()
	svc.Webclient.Server = svc.WebServer
	svc.WebApi = NewWebApi()
	svc.WebApi.Server = svc.WebServer
	svc.WebApi.svc = svc
	svc.Agent = &SvcAgent{Svc: svc}
	return svc
}

func (r *Unit) Init() error {
	r.Webclient.Init()
	r.WebApi.Init()
	if err := r.WebServer.Start(); err != nil {
		fmt.Printf("Failed to start web server: %v\n", err)
		return err
	}

	return nil
}

func (r *Unit) Run() int {

	// Initialize othe components here

	// Initialize other components here
	// Start the application
	return 0
}
