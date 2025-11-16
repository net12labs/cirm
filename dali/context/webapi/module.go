package webapi

import (
	apiwebserver "github.com/net12labs/cirm/dali/api-web-server"
	"github.com/net12labs/cirm/dali/context/cmd"
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
