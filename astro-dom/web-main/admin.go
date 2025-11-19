package domain

import (
	webclient "github.com/net12labs/cirm/astro-dom/web-admin/site-client-web"
	"github.com/net12labs/cirm/dolly/context/cmd"
	webserver "github.com/net12labs/cirm/dolly/web-server-client"
)

// this is a management interface for the domain service just like managing a vm from the outside
// we view data, statistic, logs etc., add devices, services, artefacts and so on

type domAdmin struct {
	WebClient *webclient.WebClient
	WebServer *webserver.Server
	OnExecute func(cmd *cmd.Cmd)
}

func NewDomAdmin() *domAdmin {
	ad := &domAdmin{
		WebClient: webclient.NewWebClient(),
		WebServer: webserver.NewServer(),
	}
	ad.WebClient.Server = ad.WebServer
	return ad
}

func (da *domAdmin) Init() error {
	return da.WebClient.Init()
}

func (da *domAdmin) Execute(cmd *cmd.Cmd) {}
