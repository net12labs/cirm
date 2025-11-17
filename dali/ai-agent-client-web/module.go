package aiclient

import (
	domain "github.com/net12labs/cirm/dali/domain/web-client"
	clientwebserver "github.com/net12labs/cirm/dali/web-server-client"
)

type Client struct {
	Server *clientwebserver.Server
	Domain *domain.Domain
}

func NewClient() *Client {
	return &Client{
		Domain: domain.NewDomain(),
	}
}

var NewServer = clientwebserver.NewServer

type Request = clientwebserver.Request
type Response = clientwebserver.Response
