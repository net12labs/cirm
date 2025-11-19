package webagentweb

import (
	client "github.com/net12labs/cirm/dali/client-web/agent"
)

type WebAgent struct {
	*client.ClientWeb
}

func New() *WebAgent {
	cli := &WebAgent{ClientWeb: client.NewClient()}
	cli.Domain.Path = "/agent/site/web"
	return cli
}

func (wc *WebAgent) Init() error {
	wc.Server.AddRoute(wc.Domain.Path, func(req *client.Request) error {
		return req.WriteResponse("Not implemented yet")
	})
	return nil
}
