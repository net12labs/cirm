package user

import (
	"github.com/net12labs/cirm/dali/context/service"
	webserver "github.com/net12labs/cirm/dali/web-server"

	webclient "github.com/net12labs/cirm/client-web/user"
)

type Unit struct {
	*service.Service
	WebServer   *webserver.WebServer
	Webclient   *webclient.WebClient
	Agent       *SvcAgent
	ExitMessage string
	WebApi      *WebApi
}

func NewUnit() *Unit {
	svc := &Unit{}
	svc.Service = service.NewService()
	svc.Webclient = webclient.NewWebClient()
	svc.WebApi = NewWebApi()
	svc.WebApi.svc = svc
	svc.Agent = &SvcAgent{Svc: svc}
	return svc
}

func (r *Unit) Init() error {
	r.Webclient.Server = r.WebServer
	r.Webclient.Init()
	r.WebApi.Init()
	return nil
}

func (r *Unit) Run() int {

	return 0
}
