package admin

import webserver "github.com/net12labs/cirm/dali/web-server"

type WebApi struct {
	svc    *AdminUnit
	Server *webserver.WebServer
	// WebApi fields here
}

func NewWebApi() *WebApi {
	return &WebApi{}
}

func (api *WebApi) Init() {
	// Initialize WebApi here
}
