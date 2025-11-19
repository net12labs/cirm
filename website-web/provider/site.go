package websiteweb

import (
	client "github.com/net12labs/cirm/dali/client-web/website"
)

type WebSite struct {
	*client.ClientWeb
}

func New() *WebSite {
	cli := &WebSite{ClientWeb: client.NewClient()}
	cli.Domain.Path = "/site/provider/web"
	return cli
}

func (wc *WebSite) Init() error {
	wc.Server.AddRoute(wc.Domain.Path, func(req *client.Request) error {
		return req.WriteResponse("Not implemented yet")
	})
	return nil
}
