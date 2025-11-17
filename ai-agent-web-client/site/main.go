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
	cli := &WebAiAgentClient{Client: client.NewClient()}
	cli.Domain.Path = "/site/ai-agent"
	return cli
}

func (wc *WebAiAgentClient) Init() error {
	wc.Server.AddRoute(wc.Domain.Path+"/ai-assistant", func(req *client.Request) error {
		// Serve the main page
		data, err := content.ReadFile("web/ai-assistant/index.html")
		if err != nil {
			return req.WriteResponse404()
		}
		// Wrap the HTML response with the domain information
		data = wc.Domain.WrapHTML(data, "ai-assistant", 98498479)
		return req.WriteResponseHTML(data)
	})
	return nil
}
