package client

import (
	clientwebserver "github.com/net12labs/cirm/dali/web-server/page"
)

type Client struct {
	Server      *clientwebserver.Server
	PageRequest func(req *Request)
}

func NewClient() *Client {
	return &Client{}
}

var NewServer = clientwebserver.NewServer

type Request = clientwebserver.Request
type Response = clientwebserver.Response
