package webapi

import (
	"github.com/net12labs/cirm/dolly/context/cmd"
	wpd "github.com/net12labs/cirm/dolly/domain/web-api"
	apiwebserver "github.com/net12labs/cirm/dolly/web-server-api"
)

type Request = apiwebserver.Request
type Response = apiwebserver.Response

type WebApi struct {
	Server  *apiwebserver.Server
	Execute func(cmd *cmd.Cmd)
	// WebApi fields here
}

func NewWebApi() *WebApi {
	api := &WebApi{}
	api.Domain = wpd.NewDomain()
	// Initialize WebApi fields here
	return api
}
