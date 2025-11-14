package user

import (
	"fmt"

	webserver "github.com/net12labs/cirm/dali/web-server"
	"github.com/net12labs/cirm/dali/work/service"

	webclient "github.com/net12labs/cirm/svc-client-web"
)

// Possible runmodes are; web, cli

type UserUnit struct {
	*service.ServiceUnit
	WebServer *webserver.WebServer

	// Other root fields here
	Webclient   *webclient.WebClient
	Agent       *SvcAgent
	ExitMessage string
	WebApi      *WebApi
}

func NewUserUnit() *UserUnit {
	svc := &UserUnit{}
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

func (r *UserUnit) Init() error {
	r.Webclient.Init()
	r.WebApi.Init()
	if err := r.WebServer.Start(); err != nil {
		fmt.Printf("Failed to start web server: %v\n", err)
		return err
	}

	return nil
}

func (r *UserUnit) Run() int {

	// Initialize othe components here

	// Initialize other components here
	// Start the application
	return 0
}
