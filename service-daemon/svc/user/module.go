package user

import (
	"github.com/net12labs/cirm/dali/context/service"

	webapi "github.com/net12labs/cirm/api-web/user"
	webclient "github.com/net12labs/cirm/client-web/user"
)

type Unit struct {
	*service.Service
	WebClient *webclient.WebClient
	Agent     *SvcAgent
	WebApi    *webapi.WebApi
}

func NewUnit() *Unit {
	svc := &Unit{}
	svc.Service = service.NewService()
	svc.WebClient = webclient.NewWebClient()
	svc.WebApi = webapi.NewWebApi()
	svc.Agent = &SvcAgent{Svc: svc}
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
