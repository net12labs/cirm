package webapi

import (
	"github.com/net12labs/cirm/dali/context/cmd"
	apiwebserver "github.com/net12labs/cirm/dali/web-server-api"
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
	// Initialize WebApi fields here
	return api
}
