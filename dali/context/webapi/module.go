package webapi

import (
	webserver "github.com/net12labs/cirm/dali/web-server"
)

type Request = webserver.Request
type Response = webserver.Response

type WebApi struct {
	Server *webserver.WebServer
	// WebApi fields here
}

func NewWebApi() *WebApi {
	api := &WebApi{}
	// Initialize WebApi fields here
	return api
}
