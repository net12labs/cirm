package provider

import (
	"github.com/net12labs/cirm/dali/context/service"

	webclient "github.com/net12labs/cirm/client-web/provider"
)

// Possible runmodes are; web, cli

type Unit struct {
	*service.SubService
	Service   *service.SubService
	Webclient *webclient.ProviderClient
	Agent     *SvcAgent
	WebApi    *WebApi
}

func NewUnit() *Unit {
	svc := &Unit{}
	svc.SubService = service.NewSubService()
	svc.Service = svc.SubService
	svc.Webclient = webclient.NewWebClient()
	svc.WebApi = NewWebApi()
	svc.WebApi.svc = svc
	svc.Agent = &SvcAgent{Svc: svc}
	return svc
}

func (r *Unit) Init() error {
	r.Webclient.Server = r.WebServer
	r.Webclient.Init()
	r.WebApi.Init()

	return nil
}

func (r *Unit) Run() int {

	// Initialize othe components here

	// Initialize other components here
	// Start the application
	return 0
}
