package client

import (
	webclientdomain "github.com/net12labs/cirm/dolly/domain/web-client"
	clientwebserver "github.com/net12labs/cirm/dolly/web-server-client"
)

type Client struct {
	Server *clientwebserver.Server
	Domain *webclientdomain.Domain
}

func NewClient() *Client {
	return &Client{Domain: webclientdomain.NewDomain()}
}

var NewServer = clientwebserver.NewServer

type Request = clientwebserver.Request
type Response = clientwebserver.Response
