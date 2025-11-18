package client

import (
	webapi "github.com/net12labs/cirm/dali/context/webapi"
)

type Client struct {
	*webapi.WebApi
	ApiRequest func(req *webapi.Request)
}

func NewClient() *Client {
	return &Client{WebApi: webapi.NewWebApi()}
}

type Request = webapi.Request
type Response = webapi.Response
