package webagentweb

import (
	dom "github.com/net12labs/cirm/agent-web/platform/domain"
	client "github.com/net12labs/cirm/dali/client-web/agent"
)

type WebAgent struct {
	*client.ClientWeb
}

func New() *WebAgent {
	cli := &WebAgent{ClientWeb: client.NewClient()}
	return cli
}

func (wc *WebAgent) Init() error {
	wc.Server.AddRoute(dom.Domain().Path(), func(req *client.Request) error {
		return req.WriteResponse("Not implemented yet")
	})
	return nil
}
