package site

import (
	"fmt"

	"github.com/net12labs/cirm/dali/context/cmd"
	domain_context "github.com/net12labs/cirm/dali/domain/context"
	"github.com/net12labs/cirm/dali/rtm"
	apiwebserver "github.com/net12labs/cirm/dali/web-server-api"
	clientwebserver "github.com/net12labs/cirm/dali/web-server-client"
	"github.com/net12labs/cirm/service-daemon/site/admin"
	"github.com/net12labs/cirm/service-daemon/site/platform"
	"github.com/net12labs/cirm/service-daemon/site/provider"
	"github.com/net12labs/cirm/service-daemon/site/root"
	"github.com/net12labs/cirm/service-daemon/site/user"
)

type Site struct {
	*domain_context.Domain
	User         *user.Unit
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
		User:     user.NewUnit(),
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

	sv.Root.WebClient.Server = sv.ClientServer
	sv.Provider.WebClient.Server = sv.ClientServer
	sv.Platform.WebClient.Server = sv.ClientServer
	sv.Admin.WebClient.Server = sv.ClientServer
	sv.User.WebClient.Server = sv.ClientServer

	sv.Root.WebAgentClient.Server = sv.ClientServer
	sv.Platform.WebAgentClient.Server = sv.ClientServer
	sv.Provider.WebAgentClient.Server = sv.ClientServer
	sv.Admin.WebAgentClient.Server = sv.ClientServer
	sv.User.WebAgentClient.Server = sv.ClientServer

	sv.Root.WebApi.Server = sv.ApiServer
	sv.Root.Execute = func(cmd *cmd.Cmd) {
		fmt.Println("Executing command via WebApi:", cmd)

		// User login should be handled right at the API edge
		// user create - that needs to be done by root user

		if cmd.Target == "user.login" {
			cmd.ExitCode = 0
			cmd.Result = map[string]any{"message": "Login successful", "token": "loonabalooona"}
			return
		}
		sv.Execute(cmd)
	}

	sv.Provider.WebApi.Server = sv.ApiServer
	sv.Platform.WebApi.Server = sv.ApiServer
	sv.Admin.WebApi.Server = sv.ApiServer
	sv.User.WebApi.Server = sv.ApiServer

	sv.Root.WebAgentApi.Server = sv.ApiServer
	sv.Platform.WebAgentApi.Server = sv.ApiServer
	sv.Provider.WebAgentApi.Server = sv.ApiServer
	sv.Admin.WebAgentApi.Server = sv.ApiServer
	sv.User.WebAgentApi.Server = sv.ApiServer

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

	if err := s.User.Init(); err != nil {
		rtm.Runtime.ExitErr(1, err)
	}

	return nil

}
