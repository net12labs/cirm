package provider

import (
	"fmt"

	"github.com/net12labs/cirm/dali/context/cmd"
	domain_context "github.com/net12labs/cirm/dali/domain/context"
	webclient "github.com/net12labs/cirm/site-client-web/provider"
	webapi "github.com/net12labs/cirm/site-web-api/provider"
	website "github.com/net12labs/cirm/site-web/provider"

	webagentclient "github.com/net12labs/cirm/agent-client-web/provider"
	webagentapi "github.com/net12labs/cirm/agent-web-api/provider"
	webagent "github.com/net12labs/cirm/agent-web/provider"

	aiagentwebapi "github.com/net12labs/cirm/ai-agent-web-api/provider"
	aiagentwebclient "github.com/net12labs/cirm/ai-agent-web-client/provider"
	webaiagent "github.com/net12labs/cirm/ai-agent-web/provider"
)

// Possible runmodes are; web, cli

type Unit struct {
	*domain_context.SubDomain
	Domain        *domain_context.SubDomain
	WebSiteClient *webclient.ProviderClient
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

	r.WebSiteApi.Execute = func(cmd *cmd.Cmd) {
		fmt.Println("Executing command via Platform WebSiteApi:", cmd)
		r.WebSite.OnExecute(cmd)
	}
	r.WebSite.Execute = func(cmd *cmd.Cmd) {
		fmt.Println("Executing command via Platform WebSite:", cmd)
		r.OnExecute(cmd)
	}

	r.WebAgentApi.Execute = func(cmd *cmd.Cmd) {
		fmt.Println("Executing command via Platform WebAgentApi:", cmd)
		r.WebAgent.OnExecute(cmd)
	}
	r.WebAgent.Execute = func(cmd *cmd.Cmd) {
		fmt.Println("Executing command via Platform WebAgent:", cmd)
		r.OnExecute(cmd)
	}

	r.WebAiAgentApi.Execute = func(cmd *cmd.Cmd) {
		fmt.Println("Executing command via Platform WebAiAgentApi:", cmd)
		r.WebAiAgent.OnExecute(cmd)
	}
	r.WebAiAgent.Execute = func(cmd *cmd.Cmd) {
		fmt.Println("Executing command via Platform WebAiAgent:", cmd)
		r.OnExecute(cmd)
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

	return 0
}
