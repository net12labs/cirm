package aiagentwebclient

import (
	"embed"

	dom "github.com/net12labs/cirm/astro-site/ai-agent-web-page/site/domain"
	client "github.com/net12labs/cirm/dali/client-page/ai-agent"
)

//go:embed web/*
var content embed.FS

type WebAiAgentClient struct {
	*client.Client
}

func NewClient() *WebAiAgentClient {
	cli := &WebAiAgentClient{Client: client.NewClient()}
	return cli
}

func (wc *WebAiAgentClient) Init() error {

	wc.Server.AddRoute(dom.Domain().MakePath("ai-assistant"), func(req *client.Request) error {
		// Serve the main page
		data, err := content.ReadFile("web/ai-assistant/index.html")
		if err != nil {
			return req.WriteResponse404()
		}
		return req.WriteResponseHTML(data)
	})
	return nil
}
