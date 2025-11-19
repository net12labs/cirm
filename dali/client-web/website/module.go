package clientwebwebsite

import (
	webclientdomain "github.com/net12labs/cirm/dali/domain/web-client"
	server "github.com/net12labs/cirm/dali/web-server-web"
)

type ClientWeb struct {
	Server     *server.Server
	Domain     *webclientdomain.Domain
	WebRequest func(req *Request)
}

func NewClient() *ClientWeb {
	return &ClientWeb{
		Domain: webclientdomain.NewDomain(),
	}
}

var NewServer = server.NewServer

type Request = server.Request
type Response = server.Response
