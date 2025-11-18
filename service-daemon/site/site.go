package site

import (
	"fmt"

	domain_context "github.com/net12labs/cirm/dali/domain/context"
	"github.com/net12labs/cirm/dali/rtm"
	apiwebserver "github.com/net12labs/cirm/dali/web-server-api"
	clientwebserver "github.com/net12labs/cirm/dali/web-server-client"
	webserver "github.com/net12labs/cirm/mali/web-server"
	"github.com/net12labs/cirm/service-daemon/site/admin"
	"github.com/net12labs/cirm/service-daemon/site/consumer"
	"github.com/net12labs/cirm/service-daemon/site/platform"
	"github.com/net12labs/cirm/service-daemon/site/provider"
	"github.com/net12labs/cirm/service-daemon/site/site"
)

type Site struct {
	*domain_context.Domain
	Consumer     *consumer.Unit
	Provider     *provider.Unit
	Platform     *platform.Unit
	Site         *site.Unit
	Admin        *admin.Unit
	ApiServer    *apiwebserver.Server
	ClientServer *clientwebserver.Server
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

func (sv *Site) Init() {

	sv.ApiServer = apiwebserver.NewServer()
	sv.ClientServer = clientwebserver.NewServer()

	sv.ApiServer.WebServer = sv.WebServer
	sv.ClientServer.WebServer = sv.WebServer

	sv.Site.WebSiteClient.Server = sv.ClientServer
	sv.Provider.WebSiteClient.Server = sv.ClientServer
	sv.Platform.WebSiteClient.Server = sv.ClientServer
	sv.Admin.WebSiteClient.Server = sv.ClientServer
	sv.Consumer.WebSiteClient.Server = sv.ClientServer

	sv.Site.WebAgentClient.Server = sv.ClientServer
	sv.Platform.WebAgentClient.Server = sv.ClientServer
	sv.Provider.WebAgentClient.Server = sv.ClientServer
	sv.Admin.WebAgentClient.Server = sv.ClientServer
	sv.Consumer.WebAgentClient.Server = sv.ClientServer

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

	sv.Provider.WebAiAgentClient.Server = sv.ClientServer
	sv.Admin.WebAiAgentClient.Server = sv.ClientServer
	sv.Site.WebAiAgentClient.Server = sv.ClientServer
	sv.Consumer.WebAiAgentClient.Server = sv.ClientServer
	sv.Platform.WebAiAgentClient.Server = sv.ClientServer

	// these are API request routers

	sv.Provider.ApiRequest = func(req *sitewebapi.Request) {
		fmt.Println("Executing command via WebApi:", req)
		sv.OnApiRequest(req.Request)
	}

	sv.Admin.ApiRequest = func(req *sitewebapi.Request) {
		fmt.Println("Executing command via WebApi:", req)
		sv.OnApiRequest(req.Request)
	}

	sv.Site.ApiRequest = func(req *sitewebapi.Request) {
		fmt.Println("Executing command via WebApi:", req)
		sv.OnApiRequest(req.Request)
	}

	sv.Consumer.ApiRequest = func(req *sitewebapi.Request) {
		fmt.Println("Executing command via WebApi:", req)
		sv.OnApiRequest(req.Request)
	}

	sv.Platform.ApiRequest = func(req *sitewebapi.Request) {
		fmt.Println("Executing command via WebApi:", req)
		sv.OnApiRequest(req.Request)
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
func (s *Site) OnApiRequest(req *apiwebserver.Request) {
	fmt.Println("Executing command via Site ApiRequest:", req)
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
