package clientwebagent

import (
	server "github.com/net12labs/cirm/dali/web-server/web"
)

type ClientWeb struct {
	Server     *server.Server
	WebRequest func(req *Request)
}

func NewClient() *ClientWeb {
	return &ClientWeb{}
}

var NewServer = server.NewServer

type Request = server.Request
type Response = server.Response
