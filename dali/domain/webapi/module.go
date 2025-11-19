package webapi

import (
	apiwebserver "github.com/net12labs/cirm/dali/web-server-api"
)

type Request = apiwebserver.Request
type Response = apiwebserver.Response

type WebApi struct {
	Server *apiwebserver.Server
	// WebApi fields here
}

func NewWebApi() *WebApi {
	api := &WebApi{}
	return api
}
