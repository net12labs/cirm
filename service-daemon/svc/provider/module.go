package provider

import (
	"github.com/net12labs/cirm/dali/context/service"

	webapi "github.com/net12labs/cirm/api-web/provider"
	webclient "github.com/net12labs/cirm/client-web/provider"
)

// Possible runmodes are; web, cli

type Unit struct {
	*service.SubService
	Service   *service.SubService
	WebClient *webclient.ProviderClient
	Agent     *SvcAgent
	WebApi    *webapi.WebApi
}

func NewUnit() *Unit {
	svc := &Unit{}
	svc.SubService = service.NewSubService()
	svc.Service = svc.SubService
	svc.WebClient = webclient.NewWebClient()
	svc.WebApi = webapi.NewWebApi()
	svc.Agent = &SvcAgent{Svc: svc}
	return svc
}

func (r *Unit) Init() error {
	r.WebClient.Init()
	r.WebApi.Init()

	return nil
}

func (r *Unit) Run() int {

	// Initialize othe components here

	// Initialize other components here
	// Start the application
	return 0
}
