package webaiagentweb

import (
	client "github.com/net12labs/cirm/dali/client-web/ai-agent"
)

type WebAiAgent struct {
	*client.ClientWeb
}

func New() *WebAiAgent {
	cli := &WebAiAgent{ClientWeb: client.NewClient()}
	cli.Domain.Path = "/ai-agent/platform/web"
	return cli
}

func (wc *WebAiAgent) Init() error {
	wc.Server.AddRoute(wc.Domain.Path, func(req *client.Request) error {
		return req.WriteResponse("Not implemented yet")
	})
	return nil
}
