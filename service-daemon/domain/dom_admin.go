package domain

import (
	dom_webserver "github.com/net12labs/cirm/dolly/web-server-client"
	webclient "github.com/net12labs/cirm/domain-web/site-client-web"
)

// this is a management interface for the domain service just like managing a vm from the outside
// we view data, statistic, logs etc., add devices, services, artefacts and so on

type domAdmin struct {
	WebClient *webclient.WebClient
	WebServer *dom_webserver.Server
}

func NewDomAdmin() *domAdmin {
	return &domAdmin{
		WebClient: webclient.NewWebClient(),
		WebServer: dom_webserver.NewServer(),
	}
}

func (da *domAdmin) Init() error {
	return da.WebClient.Init()
}
