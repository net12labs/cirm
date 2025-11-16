package aiagentwebclient

import (
	"embed"

	client "github.com/net12labs/cirm/dali/ai-agent-client-web"
)

//go:embed web/*
var content embed.FS

type WebAiAgentClient struct {
	*client.Client
}

func NewClient() *WebAiAgentClient {
	return &WebAiAgentClient{Client: client.NewClient()}
}

func (wc *WebAiAgentClient) Init() error {
	wc.Server.AddRoute("/admin/ai-agent", func(req *client.Request) error {
		// Serve the main page
		data, err := content.ReadFile("web/index.html")
		if err != nil {
			return req.WriteResponse404()
		}
		return req.WriteResponseHTML(data)
	})
	return nil
}
