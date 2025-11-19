package site

import (
	"fmt"

	domain_context "github.com/net12labs/cirm/dali/domain/context"
	apiwebserver "github.com/net12labs/cirm/dali/web-server-api"
	clientwebserver "github.com/net12labs/cirm/dali/web-server-client"
	web_server "github.com/net12labs/cirm/dali/web-server-web"
	"github.com/net12labs/cirm/ops/rtm"

	webserver "github.com/net12labs/cirm/mali/web-server"
	"github.com/net12labs/cirm/service-daemon/site/admin"
	"github.com/net12labs/cirm/service-daemon/site/consumer"
	"github.com/net12labs/cirm/service-daemon/site/platform"
	"github.com/net12labs/cirm/service-daemon/site/provider"
	"github.com/net12labs/cirm/service-daemon/site/site"
)

type Site struct {
	*domain_context.Domain
	Consumer   *consumer.Unit
	Provider   *provider.Unit
	Platform   *platform.Unit
	Site       *site.Unit
	Admin      *admin.Unit
	ApiServer  *apiwebserver.Server
	PageServer *clientwebserver.Server
	WebxServer *web_server.Server
}

// so the actual web server may need to be started even higher

func NewSite() *Site {
	sv := &Site{
		Domain:   domain_context.NewDomain(),
		Consumer: consumer.NewUnit(),
		Admin:    admin.NewUnit(),
		Provider: provider.NewUnit(),
		Site:     site.NewUnit(),
		Platform: platform.NewUnit(),
	}

	return sv
}

func (sv *Site) initMappings() {

	sv.ApiServer = apiwebserver.NewServer()
	sv.PageServer = clientwebserver.NewServer()
	sv.WebxServer = web_server.NewServer()

	sv.ApiServer.WebServer = sv.WebServer
	sv.PageServer.WebServer = sv.WebServer
	sv.WebxServer.WebServer = sv.WebServer

	sv.Site.WebAgent.Server = sv.WebxServer
	sv.Provider.WebAgent.Server = sv.WebxServer
	sv.Platform.WebAgent.Server = sv.WebxServer
	sv.Admin.WebAgent.Server = sv.WebxServer
	sv.Consumer.WebAgent.Server = sv.WebxServer

	sv.Site.WebAiAgent.Server = sv.WebxServer
	sv.Provider.WebAiAgent.Server = sv.WebxServer
	sv.Platform.WebAiAgent.Server = sv.WebxServer
	sv.Admin.WebAiAgent.Server = sv.WebxServer
	sv.Consumer.WebAiAgent.Server = sv.WebxServer

	sv.Site.WebSite.Server = sv.WebxServer
	sv.Provider.WebSite.Server = sv.WebxServer
	sv.Platform.WebSite.Server = sv.WebxServer
	sv.Admin.WebSite.Server = sv.WebxServer
	sv.Consumer.WebSite.Server = sv.WebxServer

	// these are Client web servers

	sv.Site.WebSiteClient.Server = sv.PageServer
	sv.Provider.WebSiteClient.Server = sv.PageServer
	sv.Platform.WebSiteClient.Server = sv.PageServer
	sv.Admin.WebSiteClient.Server = sv.PageServer
	sv.Consumer.WebSiteClient.Server = sv.PageServer

	sv.Site.WebAgentClient.Server = sv.PageServer
	sv.Platform.WebAgentClient.Server = sv.PageServer
	sv.Provider.WebAgentClient.Server = sv.PageServer
	sv.Admin.WebAgentClient.Server = sv.PageServer
	sv.Consumer.WebAgentClient.Server = sv.PageServer

	sv.Provider.WebAiAgentClient.Server = sv.PageServer
	sv.Admin.WebAiAgentClient.Server = sv.PageServer
	sv.Site.WebAiAgentClient.Server = sv.PageServer
	sv.Consumer.WebAiAgentClient.Server = sv.PageServer
	sv.Platform.WebAiAgentClient.Server = sv.PageServer

	sv.Site.WebSiteApi.Server = sv.ApiServer
	sv.Provider.WebSiteApi.Server = sv.ApiServer
	sv.Platform.WebSiteApi.Server = sv.ApiServer
	sv.Admin.WebSiteApi.Server = sv.ApiServer
	sv.Consumer.WebSiteApi.Server = sv.ApiServer

	sv.Site.WebAgentApi.Server = sv.ApiServer
	sv.Platform.WebAgentApi.Server = sv.ApiServer
	sv.Provider.WebAgentApi.Server = sv.ApiServer
	sv.Admin.WebAgentApi.Server = sv.ApiServer
	sv.Consumer.WebAgentApi.Server = sv.ApiServer

	sv.Provider.WebAiAgentApi.Server = sv.ApiServer
	sv.Admin.WebAiAgentApi.Server = sv.ApiServer
	sv.Site.WebAiAgentApi.Server = sv.ApiServer
	sv.Consumer.WebAiAgentApi.Server = sv.ApiServer
	sv.Platform.WebAiAgentApi.Server = sv.ApiServer

}

func (sv *Site) Init() {

	sv.initMappings()
	// these are Web request routers
	sv.Provider.WebRequest = func(req *webserver.Request) {
		fmt.Println("Executing command via WebApi:", req)
		sv.OnWebRequest(req)
	}

	sv.Admin.WebRequest = func(req *webserver.Request) {
		fmt.Println("Executing command via WebApi:", req)
		sv.OnWebRequest(req)
	}

	sv.Site.WebRequest = func(req *webserver.Request) {
		fmt.Println("Executing command via WebApi:", req)
		sv.OnWebRequest(req)
	}

	sv.Consumer.WebRequest = func(req *webserver.Request) {
		fmt.Println("Executing command via WebApi:", req)
		sv.OnWebRequest(req)
	}

	sv.Platform.WebRequest = func(req *webserver.Request) {
		fmt.Println("Executing command via WebApi:", req)
		sv.OnWebRequest(req)
	}

	// these are API request routers
	sv.Provider.ApiRequest = func(req *webserver.Request) {
		fmt.Println("Executing command via WebApi:", req)
		sv.OnApiRequest(req)
	}

	sv.Admin.ApiRequest = func(req *webserver.Request) {
		fmt.Println("Executing command via WebApi:", req)
		sv.OnApiRequest(req)
	}

	sv.Site.ApiRequest = func(req *webserver.Request) {
		fmt.Println("Executing command via WebApi:", req)
		sv.OnApiRequest(req)
	}

	sv.Consumer.ApiRequest = func(req *webserver.Request) {
		fmt.Println("Executing command via WebApi:", req)
		sv.OnApiRequest(req)
	}

	sv.Platform.ApiRequest = func(req *webserver.Request) {
		fmt.Println("Executing command via WebApi:", req)
		sv.OnApiRequest(req)
	}

	// Page request handlers

	sv.Provider.PageRequest = func(req *webserver.Request) {
		fmt.Println("Getting html page at the root", req)
		sv.OnPageRequest(req)
	}

	sv.Consumer.PageRequest = func(req *webserver.Request) {
		fmt.Println("Getting html page at the root", req)
		sv.OnPageRequest(req)
	}

	sv.Platform.PageRequest = func(req *webserver.Request) {
		fmt.Println("Getting html page at the root", req)
		sv.OnPageRequest(req)
	}
	sv.Admin.PageRequest = func(req *webserver.Request) {
		fmt.Println("Getting html page at the root", req)
		sv.OnPageRequest(req)
	}
	sv.Site.PageRequest = func(req *webserver.Request) {
		fmt.Println("Getting html page at the root", req)
		sv.OnPageRequest(req)
	}

}

func (s *Site) OnPageRequest(req *webserver.Request) {
	fmt.Println("Executing command via Site PageRequest:", req)
}
func (s *Site) OnApiRequest(req *webserver.Request) {
	fmt.Println("Executing command via Site ApiRequest:", req)
}

func (s *Site) OnWebRequest(req *webserver.Request) {
	fmt.Println("Executing command via Site WebRequest:", req)
}

func (s *Site) Start() error {

	if err := s.startSubdomains(); err != nil {
		return err
	}
	return nil
}

func (s *Site) startSubdomains() error {

	if err := s.Site.Init(); err != nil {
		rtm.Runtime.ExitErr(1, err)
	}

	if err := s.Admin.Init(); err != nil {
		rtm.Runtime.ExitErr(1, err)
	}

	if err := s.Platform.Init(); err != nil {
		rtm.Runtime.ExitErr(1, err)
	}

	if err := s.Provider.Init(); err != nil {
		rtm.Runtime.ExitErr(1, err)
	}

	if err := s.Consumer.Init(); err != nil {
		rtm.Runtime.ExitErr(1, err)
	}

	return nil

}
