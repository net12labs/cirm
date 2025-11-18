package site

import (
	"fmt"

	"github.com/net12labs/cirm/dali/context/cmd"
	domain_context "github.com/net12labs/cirm/dali/domain/context"
	webclient "github.com/net12labs/cirm/site-client-web/site"
	webapi "github.com/net12labs/cirm/site-web-api/site"
	website "github.com/net12labs/cirm/site-web/site"

	sitewebapi "github.com/net12labs/cirm/dali/site-client-api"

	webagentclient "github.com/net12labs/cirm/agent-client-web/site"
	webagentapi "github.com/net12labs/cirm/agent-web-api/site"
	webagent "github.com/net12labs/cirm/agent-web/site"

	aiagentwebapi "github.com/net12labs/cirm/ai-agent-web-api/site"
	aiagentwebclient "github.com/net12labs/cirm/ai-agent-web-client/site"
	webaiagent "github.com/net12labs/cirm/ai-agent-web/site"
)

// Possible runmodes are; web, cli

type Unit struct {
	*domain_context.SubDomain
	Domain        *domain_context.SubDomain
	WebSiteClient *webclient.WebClient
	WebSiteApi    *webapi.WebApi
	WebSite       *website.Site

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
	svc.WebSite = website.NewSite()

	svc.WebAgent = webagent.NewAgent()
	svc.WebAgentApi = webagentapi.NewWebApi()
	svc.WebAgentClient = webagentclient.NewClient()

	svc.WebAiAgent = webaiagent.NewAiAgent()
	svc.WebAiAgentApi = aiagentwebapi.NewWebApi()
	svc.WebAiAgentClient = aiagentwebclient.NewClient()
	return svc
}

func (r *Unit) Init() error {

	// todo: so these need to be All converted to the ApiRequest

	r.WebSiteApi.ApiRequest = func(req *sitewebapi.Request) {
		fmt.Println("Executing command via Platform WebSite:", req)
		r.OnApiRequest(req.Request)
	}

	r.WebAgentApi.Execute = func(cmd *cmd.Cmd) {
		fmt.Println("Executing command via Platform WebAgentApi:", cmd)
		r.WebAgent.OnExecute(cmd)
	}
	r.WebAgent.Execute = func(cmd *cmd.Cmd) {
		fmt.Println("Executing command via Platform WebAgent:", cmd)
	}

	r.WebAiAgentApi.Execute = func(cmd *cmd.Cmd) {
		fmt.Println("Executing command via Platform WebAiAgentApi:", cmd)
		r.WebAiAgent.OnExecute(cmd)
	}
	r.WebAiAgent.Execute = func(cmd *cmd.Cmd) {
		fmt.Println("Executing command via Platform WebAiAgent:", cmd)
	}

	r.WebSite.Init()
	r.WebAgent.Init()
	r.WebAiAgent.Init()

	r.WebSiteApi.Init()
	r.WebAgentApi.Init()
	r.WebAiAgentApi.Init()

	r.WebSiteClient.Init()
	r.WebAgentClient.Init()
	r.WebAiAgentClient.Init()

	return nil
}

func (r *Unit) Run() int {

	// Initialize other components here
	// Start the application
	return 0
}
