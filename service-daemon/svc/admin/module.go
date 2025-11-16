package admin

import (
	"github.com/net12labs/cirm/dali/context/service"

	webapi "github.com/net12labs/cirm/api-web/admin"
	webclient "github.com/net12labs/cirm/client-web/admin"
)

// Possible runmodes are; web, cli

type Unit struct {
	*service.SubService
	Service   *service.SubService
	WebClient *webclient.WebClient
	Agent     *SvcAgent
	WebApi    *webapi.WebApi
}

func NewUnit() *Unit {
	svc := &Unit{}
	svc.SubService = service.NewSubService()
	svc.Service = svc.SubService
	svc.WebClient = webclient.NewWebClient()
	svc.WebApi = webapi.NewWebApi()
	// svc.WebApi.svc = svc
	svc.Agent = &SvcAgent{Svc: svc}

	return svc
}

func (r *Unit) Init() error {
	r.WebClient.Init()
	r.WebApi.Init()
	return nil
}

func (r *Unit) Run() int {

	// Initialize other components here
	// Start the application
	return 0
}
