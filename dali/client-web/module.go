package client

import webserver "github.com/net12labs/cirm/dali/web-server"

type Client struct {
	Server *webserver.WebServer
}

func NewClient() *Client {
	return &Client{}
}

var NewWebServer = webserver.NewWebServer
