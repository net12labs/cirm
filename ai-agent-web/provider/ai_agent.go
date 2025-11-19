package webaiagentweb

import (
	dom "github.com/net12labs/cirm/ai-agent-web/provider/domain"
	client "github.com/net12labs/cirm/dali/client-web/ai-agent"
)

type WebAiAgent struct {
	*client.ClientWeb
}

func New() *WebAiAgent {
	cli := &WebAiAgent{ClientWeb: client.NewClient()}
	return cli
}

func (wc *WebAiAgent) Init() error {
	wc.Server.AddRoute(dom.Domain().Path(), func(req *client.Request) error {
		return req.WriteResponse("Not implemented yet")
	})
	return nil
}
