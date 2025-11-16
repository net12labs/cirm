package user

import (
	"github.com/net12labs/cirm/dali/context/service"

	webapi "github.com/net12labs/cirm/api-web/user"
	webclient "github.com/net12labs/cirm/client-web/user"

	webagentclient "github.com/net12labs/cirm/agent-client-web/user"
	webagentapi "github.com/net12labs/cirm/agent-web-api/user"
	webagent "github.com/net12labs/cirm/agent-web/user"
)

type Unit struct {
	*service.Service
	WebClient      *webclient.WebClient
	WebApi         *webapi.WebApi
	WebAgent       *webagent.Agent
	WebAgentApi    *webagentapi.WebAgentApi
	WebAgentClient *webagentclient.WebAgentClient
}

func NewUnit() *Unit {
	svc := &Unit{}
	svc.Service = service.NewService()
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
	return nil
}

func (r *Unit) Run() int {

	return 0
}
