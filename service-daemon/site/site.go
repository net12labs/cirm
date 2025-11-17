package site

import (
	"fmt"

	"github.com/net12labs/cirm/dali/context/cmd"
	domain_context "github.com/net12labs/cirm/dali/domain/context"
	"github.com/net12labs/cirm/dali/rtm"
	apiwebserver "github.com/net12labs/cirm/dali/web-server-api"
	clientwebserver "github.com/net12labs/cirm/dali/web-server-client"
	"github.com/net12labs/cirm/service-daemon/site/admin"
	"github.com/net12labs/cirm/service-daemon/site/consumer"
	"github.com/net12labs/cirm/service-daemon/site/platform"
	"github.com/net12labs/cirm/service-daemon/site/provider"
	"github.com/net12labs/cirm/service-daemon/site/root"
)

type Site struct {
	*domain_context.Domain
	Consumer     *consumer.Unit
	Provider     *provider.Unit
	Platform     *platform.Unit
	Root         *root.Unit
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
		Root:     root.NewUnit(),
		Platform: platform.NewUnit(),
	}

	return sv
}

func (sv *Site) Init() {

	sv.ApiServer = apiwebserver.NewServer()
	sv.ClientServer = clientwebserver.NewServer()

	sv.ApiServer.WebServer = sv.WebServer
	sv.ClientServer.WebServer = sv.WebServer

	sv.Root.WebSiteClient.Server = sv.ClientServer
	sv.Provider.WebSiteClient.Server = sv.ClientServer
	sv.Platform.WebSiteClient.Server = sv.ClientServer
	sv.Admin.WebSiteClient.Server = sv.ClientServer
	sv.Consumer.WebSiteClient.Server = sv.ClientServer

	sv.Root.WebAgentClient.Server = sv.ClientServer
	sv.Platform.WebAgentClient.Server = sv.ClientServer
	sv.Provider.WebAgentClient.Server = sv.ClientServer
	sv.Admin.WebAgentClient.Server = sv.ClientServer
	sv.Consumer.WebAgentClient.Server = sv.ClientServer

	sv.Root.WebSiteApi.Server = sv.ApiServer
	sv.Provider.WebSiteApi.Server = sv.ApiServer
	sv.Platform.WebSiteApi.Server = sv.ApiServer
	sv.Admin.WebSiteApi.Server = sv.ApiServer
	sv.Consumer.WebSiteApi.Server = sv.ApiServer

	sv.Root.WebAgentApi.Server = sv.ApiServer
	sv.Platform.WebAgentApi.Server = sv.ApiServer
	sv.Provider.WebAgentApi.Server = sv.ApiServer
	sv.Admin.WebAgentApi.Server = sv.ApiServer
	sv.Consumer.WebAgentApi.Server = sv.ApiServer

	sv.Provider.WebAiAgentApi.Server = sv.ApiServer
	sv.Admin.WebAiAgentApi.Server = sv.ApiServer
	sv.Root.WebAiAgentApi.Server = sv.ApiServer
	sv.Consumer.WebAiAgentApi.Server = sv.ApiServer
	sv.Platform.WebAiAgentApi.Server = sv.ApiServer

	sv.Provider.WebAiAgentClient.Server = sv.ClientServer
	sv.Admin.WebAiAgentClient.Server = sv.ClientServer
	sv.Root.WebAiAgentClient.Server = sv.ClientServer
	sv.Consumer.WebAiAgentClient.Server = sv.ClientServer
	sv.Platform.WebAiAgentClient.Server = sv.ClientServer

	sv.Provider.Execute = func(cmd *cmd.Cmd) {
		fmt.Println("Executing command via WebApi:", cmd)
		sv.OnExecute(cmd)
	}

	sv.Admin.Execute = func(cmd *cmd.Cmd) {
		fmt.Println("Executing command via WebApi:", cmd)
		sv.OnExecute(cmd)
	}

	sv.Root.Execute = func(cmd *cmd.Cmd) {
		fmt.Println("Executing command via WebApi:", cmd)
		sv.OnExecute(cmd)
	}

	sv.Consumer.Execute = func(cmd *cmd.Cmd) {
		fmt.Println("Executing command via WebApi:", cmd)
		sv.OnExecute(cmd)
	}

	sv.Platform.Execute = func(cmd *cmd.Cmd) {
		fmt.Println("Executing command via WebApi:", cmd)
		sv.OnExecute(cmd)
	}

}

func (s *Site) Start() error {

	if err := s.startSubdomains(); err != nil {
		return err
	}
	return nil
}

func (s *Site) startSubdomains() error {

	if err := s.Root.Init(); err != nil {
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
