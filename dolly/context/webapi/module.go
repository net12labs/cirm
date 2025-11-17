package webapi

import (
	"github.com/net12labs/cirm/dolly/context/cmd"
	webapidomain "github.com/net12labs/cirm/dolly/domain/web-api"
	apiwebserver "github.com/net12labs/cirm/dolly/web-server-api"
)

type Request = apiwebserver.Request
type Response = apiwebserver.Response

type WebApi struct {
	Server  *apiwebserver.Server
	Execute func(cmd *cmd.Cmd)
	Domain  *webapidomain.Domain
	// WebApi fields here
}

func NewWebApi() *WebApi {
	api := &WebApi{}
	api.Domain = webapidomain.NewDomain()
	// Initialize WebApi fields here
	return api
}
