package provider

import (
	"github.com/net12labs/cirm/dali/context/service"

	webclient "github.com/net12labs/cirm/site-client-web/provider"
	webapi "github.com/net12labs/cirm/site-web-api/provider"

	webagentclient "github.com/net12labs/cirm/agent-client-web/provider"
	webagentapi "github.com/net12labs/cirm/agent-web-api/provider"
	webagent "github.com/net12labs/cirm/agent-web/provider"
)

// Possible runmodes are; web, cli

type Unit struct {
	*service.SubService
	Service        *service.SubService
	WebClient      *webclient.ProviderClient
	WebApi         *webapi.WebApi
	WebAgent       *webagent.Agent
	WebAgentApi    *webagentapi.WebAgentApi
	WebAgentClient *webagentclient.WebAgentClient
}

func NewUnit() *Unit {
	svc := &Unit{}
	svc.SubService = service.NewSubService()
	svc.Service = svc.SubService
	svc.WebClient = webclient.NewWebClient()
	svc.WebApi = webapi.NewWebApi()
	svc.WebAgent = webagent.NewAgent()
	svc.WebAgentApi = webagentapi.NewWebApi()
	svc.WebAgentClient = webagentclient.NewClient()
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
