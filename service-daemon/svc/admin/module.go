package admin

import (
	"fmt"

	webserver "github.com/net12labs/cirm/dali/web-server"
	"github.com/net12labs/cirm/dali/work/service"

	webclient "github.com/net12labs/cirm/admin-client-web"
)

// Possible runmodes are; web, cli

type AdminUnit struct {
	*service.ServiceUnit
	WebServer *webserver.WebServer
	WebClient *webclient.AdminClient
	// Other root fields here

	Agent       *SvcAgent
	ExitMessage string
	WebApi      *WebApi
}

func NewAdminUnit() *AdminUnit {
	svc := &AdminUnit{}
	svc.ServiceUnit = service.NewServiceUnit()
	svc.WebServer = webserver.NewWebServer()
	svc.WebClient = webclient.NewAdminClient()
	svc.WebClient.Server = svc.WebServer
	svc.WebApi = NewWebApi()
	svc.WebApi.Server = svc.WebServer
	svc.WebApi.svc = svc
	svc.Agent = &SvcAgent{Svc: svc}

	return svc
}

func (r *AdminUnit) Init() error {
	r.WebClient.Init()
	r.WebApi.Init()
	if err := r.WebServer.Start(); err != nil {
		fmt.Printf("Failed to start web server: %v\n", err)
		return err
	}

	return nil
}

func (r *AdminUnit) Run() int {

	// Initialize othe components here

	// Initialize other components here
	// Start the application
	return 0
}
