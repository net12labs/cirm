package user

import (
	domain_context "github.com/net12labs/cirm/dali/domain/context"
	webclient "github.com/net12labs/cirm/site-client-web/user"
	webapi "github.com/net12labs/cirm/site-web-api/user"

	webagentclient "github.com/net12labs/cirm/agent-client-web/user"
	webagentapi "github.com/net12labs/cirm/agent-web-api/user"
	webagent "github.com/net12labs/cirm/agent-web/user"

	aiagentwebapi "github.com/net12labs/cirm/ai-agent-web-api/admin"
	aiagentwebclient "github.com/net12labs/cirm/ai-agent-web-client/admin"
	webaiagent "github.com/net12labs/cirm/ai-agent-web/admin"
)

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
	svc.Domain = domain_context.NewSubDomain()
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
	return nil
}

func (r *Unit) Run() int {

	return 0
}
