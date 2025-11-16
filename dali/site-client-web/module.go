package client

import clientwebserver "github.com/net12labs/cirm/dali/web-server-client"

type Client struct {
	Server *clientwebserver.Server
}

func NewClient() *Client {
	return &Client{}
}

var NewServer = clientwebserver.NewServer

type Request = clientwebserver.Request
type Response = clientwebserver.Response
