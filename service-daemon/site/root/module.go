package root

import (
	"fmt"

	webclient "github.com/net12labs/cirm/client-web/root"
	"github.com/net12labs/cirm/dali/context/cmd"
	"github.com/net12labs/cirm/dali/context/service"
	webapi "github.com/net12labs/cirm/site-web-api/root"

	webagentclient "github.com/net12labs/cirm/agent-client-web/root"
	webagentapi "github.com/net12labs/cirm/agent-web-api/root"
	webagent "github.com/net12labs/cirm/agent-web/root"
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
	r.WebAgentClient.Init()

	r.WebAgentApi.Execute = func(cmd *cmd.Cmd) {
		fmt.Println("Executing command via WebAgentApi:", cmd)
		r.Execute(cmd)
		// Implement command execution logic here
	}

	r.WebApi.Execute = func(cmd *cmd.Cmd) {
		fmt.Println("Executing command via WebApi:", cmd)
		r.Execute(cmd)
		// Implement command execution logic here
	}

	return nil
}

func (r *Unit) Run() int {

	// Initialize other components here
	// Start the application
	return 0
}
