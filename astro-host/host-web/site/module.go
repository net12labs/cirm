package site

import (
	"fmt"

	domain_context "github.com/net12labs/cirm/dali/domain/context"

	webapi "github.com/net12labs/cirm/astro-site/web/website-web-api/site"
	webclient "github.com/net12labs/cirm/astro-site/web/website-web-page/site"
	website "github.com/net12labs/cirm/astro-site/web/website-web/site"

	webagentclient "github.com/net12labs/cirm/astro-site/web/agent-client-page/site"
	webagentapi "github.com/net12labs/cirm/astro-site/web/agent-web-api/site"
	webagent "github.com/net12labs/cirm/astro-site/web/agent-web/site"

	aiagentwebapi "github.com/net12labs/cirm/astro-site/web/ai-agent-web-api/site"
	aiagentwebclient "github.com/net12labs/cirm/astro-site/web/ai-agent-web-page/site"
	webaiagent "github.com/net12labs/cirm/astro-site/web/ai-agent-web/site"

	agent_client "github.com/net12labs/cirm/dali/client-page/agent"
	aiagent_client "github.com/net12labs/cirm/dali/client-page/ai-agent"
	website_client "github.com/net12labs/cirm/dali/client-page/website"

	agent_api "github.com/net12labs/cirm/dali/client-api/agent"
	aiagent_api "github.com/net12labs/cirm/dali/client-api/ai-agent"
	website_api "github.com/net12labs/cirm/dali/client-api/website"

	agent_web "github.com/net12labs/cirm/dali/client-web/agent"
	aiagent_web "github.com/net12labs/cirm/dali/client-web/ai-agent"
	website_web "github.com/net12labs/cirm/dali/client-web/website"
)

// Possible runmodes are; web, cli

type Unit struct {
	*domain_context.SubDomain
	Domain        *domain_context.SubDomain
	WebSiteClient *webclient.WebClient
	WebSiteApi    *webapi.WebApi
	WebSite       *website.WebSite

	WebAgent       *webagent.WebAgent
	WebAgentApi    *webagentapi.WebAgentApi
	WebAgentClient *webagentclient.WebAgentClient

	WebAiAgent       *webaiagent.WebAiAgent
	WebAiAgentApi    *aiagentwebapi.WebAiAgentApi
	WebAiAgentClient *aiagentwebclient.WebAiAgentClient
}

func NewUnit() *Unit {
	svc := &Unit{}
	svc.SubDomain = domain_context.NewSubDomain()
	svc.Domain = svc.SubDomain

	svc.WebSiteClient = webclient.NewWebClient()
	svc.WebSiteApi = webapi.NewWebApi()
	svc.WebSite = website.New()

	svc.WebAgent = webagent.New()
	svc.WebAgentApi = webagentapi.NewWebApi()
	svc.WebAgentClient = webagentclient.NewClient()

	svc.WebAiAgent = webaiagent.New()
	svc.WebAiAgentApi = aiagentwebapi.NewWebApi()
	svc.WebAiAgentClient = aiagentwebclient.NewClient()
	return svc
}

func (r *Unit) Init() error {

	r.WebSite.WebRequest = func(req *website_web.Request) {
		fmt.Println("Executing command via Admin WebSite:", req)
		r.OnWebRequest(req.Request)
	}
	r.WebAgent.WebRequest = func(req *agent_web.Request) {
		fmt.Println("Executing command via Admin WebAgent:", req)
		r.OnWebRequest(req.Request)
	}

	r.WebAiAgent.WebRequest = func(req *aiagent_web.Request) {
		fmt.Println("Executing command via Admin WebAiAgentApi:", req)
		r.OnWebRequest(req.Request)
	}

	r.WebSite.Init()
	r.WebAgent.Init()
	r.WebAiAgent.Init()

	r.WebSiteApi.ApiRequest = func(req *website_api.Request) {
		fmt.Println("Executing command via Platform WebSite:", req)
		r.OnApiRequest(req.Request)
	}

	r.WebAgentApi.ApiRequest = func(req *agent_api.Request) {
		fmt.Println("Executing command via Platform WebAgent:", req)
		r.OnApiRequest(req.Request)
	}

	r.WebAiAgentApi.ApiRequest = func(req *aiagent_api.Request) {
		fmt.Println("Executing command via Platform WebAiAgent:", req)
		r.OnApiRequest(req.Request)
	}

	r.WebSiteApi.Init()
	r.WebAgentApi.Init()
	r.WebAiAgentApi.Init()

	r.WebSiteClient.PageRequest = func(req *website_client.Request) {
		fmt.Println("Executing command via Admin WebSite:", req)
		r.OnApiRequest(req.Request)
	}

	r.WebAgentClient.PageRequest = func(req *agent_client.Request) {
		fmt.Println("Executing command via Admin WebAgent:", req)
		r.OnApiRequest(req.Request)
	}

	r.WebAiAgentClient.PageRequest = func(req *aiagent_client.Request) {
		fmt.Println("Executing command via Admin WebAiAgent:", req)
		r.OnApiRequest(req.Request)
	}

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
