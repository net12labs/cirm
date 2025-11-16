package admin

import (
	webagentclient "github.com/net12labs/cirm/agent-client-web/admin"
	webagentapi "github.com/net12labs/cirm/agent-web-api/admin"
	domain_context "github.com/net12labs/cirm/dali/domain/context"

	webagent "github.com/net12labs/cirm/agent-web/admin"
	webclient "github.com/net12labs/cirm/site-client-web/admin"
	webapi "github.com/net12labs/cirm/site-web-api/admin"

	aiagentwebapi "github.com/net12labs/cirm/ai-agent-web-api/admin"
	aiagentwebclient "github.com/net12labs/cirm/ai-agent-web-client/admin"
	webaiagent "github.com/net12labs/cirm/ai-agent-web/admin"
)

// Possible runmodes are; web, cli

type Unit struct {
	*domain_context.SubDomain
	Domain         *domain_context.SubDomain
	WebSiteClient  *webclient.WebClient
	WebSiteApi     *webapi.WebApi
	WebAgent       *webagent.Agent
	WebAgentApi    *webagentapi.WebAgentApi
	WebAgentClient *webagentclient.WebAgentClient

	WebAiAgent       *webaiagent.AiAgent
	WebAiAgentApi    *aiagentwebapi.WebAiAgentApi
	WebAiAgentClient *aiagentwebclient.WebAiAgentClient
}

func NewUnit() *Unit {
	svc := &Unit{}
	svc.SubDomain = domain_context.NewSubDomain()
	svc.Domain = svc.SubDomain
	svc.WebSiteClient = webclient.NewWebClient()
	svc.WebSiteApi = webapi.NewWebApi()
	svc.WebAgent = webagent.NewAgent()
	svc.WebAgentApi = webagentapi.NewWebApi()
	svc.WebAgentClient = webagentclient.NewClient()

	svc.WebAiAgent = webaiagent.NewAiAgent()
	svc.WebAiAgentApi = aiagentwebapi.NewWebApi()
	svc.WebAiAgentClient = aiagentwebclient.NewClient()

	return svc
}

func (r *Unit) Init() error {
	r.WebSiteClient.Init()
	r.WebSiteApi.Init()
	r.WebAgent.Init()
	r.WebAgentApi.Init()
	return nil
}

func (r *Unit) Run() int {

	// Initialize other components here
	// Start the application
	return 0
}
