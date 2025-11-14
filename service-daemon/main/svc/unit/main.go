package unit

import (
	webserver "cirm/lib/web-server"
	"cirm/lib/work/service"
	"cirm/mod/config"
	webclient "cirm/svc/web-client"
	"fmt"
)

// Possible runmodes are; web, cli

type SvcUnit struct {
	*service.ServiceUnit
	Config    config.Config
	WebServer *webserver.WebServer

	// Other root fields here
	Webclient   *webclient.WebClient
	Agent       *SvcAgent
	ExitMessage string
	WebApi      *WebApi
}

func NewSvcUnit() *SvcUnit {
	svc := &SvcUnit{}
	svc.ServiceUnit = service.NewServiceUnit()
	svc.WebServer = webserver.NewWebServer()
	svc.Webclient = webclient.NewWebClient()
	svc.Webclient.Server = svc.WebServer
	svc.WebApi = NewWebApi()
	svc.WebApi.Server = svc.WebServer
	svc.WebApi.svc = svc
	svc.Agent = &SvcAgent{Svc: svc}

	return svc
}

func (r *SvcUnit) Init() error {
	r.Webclient.Init()
	r.WebApi.Init()
	if err := r.WebServer.Start(); err != nil {
		fmt.Printf("Failed to start web server: %v\n", err)
		return err
	}

	return nil
}

func (r *SvcUnit) Run() int {

	// Initialize othe components here

	// Initialize other components here
	// Start the application
	return 0
}
