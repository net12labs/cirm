package websiteweb

import (
	client "github.com/net12labs/cirm/dali/client-web/website"
	dom "github.com/net12labs/cirm/website-web/platform/domain"
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
