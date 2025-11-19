package client

import (
	webapi "github.com/net12labs/cirm/dali/context/webapi"
)

type ClientApi struct {
	*webapi.WebApi
	ApiRequest func(req *webapi.Request)
}

func NewClient() *ClientApi {
	return &ClientApi{WebApi: webapi.NewWebApi()}
}

type Request = webapi.Request
type Response = webapi.Response
