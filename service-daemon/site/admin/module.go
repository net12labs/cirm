package admin

import (
	"github.com/net12labs/cirm/dali/context/service"

	webagentclient "github.com/net12labs/cirm/agent-client-web/admin"
	webagentapi "github.com/net12labs/cirm/agent-web-api/admin"

	webagent "github.com/net12labs/cirm/agent-web/admin"
	webclient "github.com/net12labs/cirm/site-client-web/admin"
	webapi "github.com/net12labs/cirm/site-web-api/admin"
)

// Possible runmodes are; web, cli

type Unit struct {
	*service.SubService
	Service        *service.SubService
	WebClient      *webclient.WebClient
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
	r.WebAgent.Init()
	r.WebAgentApi.Init()
	return nil
}

func (r *Unit) Run() int {

	// Initialize other components here
	// Start the application
	return 0
}
