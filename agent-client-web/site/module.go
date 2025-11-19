package webagentclient

import (
	"embed"

	client "github.com/net12labs/cirm/dali/client-page/agent"
)

//go:embed web/*
var content embed.FS

type WebAgentClient struct {
	*client.Client
}

func NewClient() *WebAgentClient {
	cli := &WebAgentClient{Client: client.NewClient()}
	cli.Domain.Path = "/site/agent"
	return cli
}

func (wc *WebAgentClient) Init() error {
	// AI Assistant partial widget
	wc.Server.AddRoute(wc.Domain.MakePath("account-create"), func(req *client.Request) error {
		data, err := content.ReadFile("web/account-create/index.html")
		if err != nil {
			return req.WriteResponse404()
		}
		data = wc.Domain.WrapHTML(data, "account-create", 98498479)
		return req.WriteResponseHTML(data)
	})

	// Login partial widget
	wc.Server.AddRoute(wc.Domain.MakePath("login"), func(req *client.Request) error {
		data, err := content.ReadFile("web/login/index.html")
		if err != nil {
			return req.WriteResponse404()
		}
		data = wc.Domain.WrapHTML(data, "login", 98498479)
		return req.WriteResponseHTML(data)
	})

	wc.Server.AddRoute(wc.Domain.MakePath("assistant"), func(req *client.Request) error {
		data, err := content.ReadFile("web/assistant/index.html")
		if err != nil {
			return req.WriteResponse404()
		}
		data = wc.Domain.WrapHTML(data, "assistant", 98498479)
		return req.WriteResponseHTML(data)
	})

	return nil
}
