package webapi

import (
	apiwebserver "github.com/net12labs/cirm/dali/api-web-server"
)

type Request = apiwebserver.ApiRequest
type Response = apiwebserver.ApiResponse

type WebApi struct {
	Server *apiwebserver.Server
	// WebApi fields here
}

func NewWebApi() *WebApi {
	api := &WebApi{}
	// Initialize WebApi fields here
	return api
}
