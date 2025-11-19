package webagentclient

import (
	"embed"

	dom "github.com/net12labs/cirm/agent-client-web/provider/domain"
	client "github.com/net12labs/cirm/dali/client-page/agent"
)

//go:embed web/*
var content embed.FS

type WebAgentClient struct {
	*client.Client
}

func (w *WebAgentClient) Domain() *dom.Dom {
	return dom.Domain()
}

func NewClient() *WebAgentClient {
	cli := &WebAgentClient{Client: client.NewClient()}
	return cli
}

func (wc *WebAgentClient) Init() error {
	wc.Server.AddRoute(dom.Domain().Path(), func(req *client.Request) error {
		data, err := content.ReadFile("web/index.html")
		if err != nil {
			return req.WriteResponse404()
		}
		return req.WriteResponseHTML(data)
	})
	return nil
}
