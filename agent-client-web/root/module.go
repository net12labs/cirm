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
	// AI Assistant partial widget
	wc.Server.AddRoute("/agent/account-create", func(req *client.Request) error {
		data, err := content.ReadFile("web/account-create/index.html")
		if err != nil {
			return req.WriteResponse404()
		}
		return req.WriteResponseHTML(data)
	})

	// Login partial widget
	wc.Server.AddRoute("/agent/login", func(req *client.Request) error {
		data, err := content.ReadFile("web/login/index.html")
		if err != nil {
			return req.WriteResponse404()
		}
		return req.WriteResponseHTML(data)
	})

	wc.Server.AddRoute("/agent/assistant", func(req *client.Request) error {
		if req.Path.Path != "/agent/assistant" {
			return req.WriteResponse404()
		}
		data, err := content.ReadFile("web/assistant/index.html")
		if err != nil {
			return req.WriteResponse404()
		}
		return req.WriteResponseHTML(data)
	})

	return nil
}
