package websiteweb

import (
	dom "github.com/net12labs/cirm/astro-site/website-web/consumer/domain"
	client "github.com/net12labs/cirm/dali/client-web/website"
)

type WebSite struct {
	*client.ClientWeb
}

func New() *WebSite {
	cli := &WebSite{ClientWeb: client.NewClient()}
	return cli
}

func (wc *WebSite) Init() error {
	wc.Server.AddRoute(dom.Domain().Path(), func(req *client.Request) error {
		return req.WriteResponse("Not implemented yet")
	})
	return nil
}
