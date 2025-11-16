package webagentclient

import (
	"embed"

	client "github.com/net12labs/cirm/dali/agent-client-web"
)

//go:embed web/*
var content embed.FS

type WebAgentClient struct {
	*client.Client
}

func NewClient() *WebAgentClient {
	return &WebAgentClient{Client: client.NewClient()}
}

func (wc *WebAgentClient) Init() error {
	wc.Server.AddRoute("/user/agent", func(req *client.Request) error {

		data, err := content.ReadFile("web/index.html")
		if err != nil {
			return req.WriteResponse404()
		}
		return req.WriteResponseHTML(data)
	})
	return nil
}
