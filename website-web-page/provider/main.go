package webclient

import (
	"embed"

	client "github.com/net12labs/cirm/dali/client-page/website"
)

//go:embed web/*
var content embed.FS

type WebClient struct {
	*client.Client
}

func NewWebClient() *WebClient {
	cl := &WebClient{Client: client.NewClient()}
	cl.Domain.Path = "/provider"
	return cl
}

func (wc *WebClient) Init() error {
	wc.Server.AddRoute(wc.Domain.Path, func(req *client.Request) error {
		data, err := content.ReadFile("web/index.html")
		if err != nil {
			return req.WriteResponse404()
		}
		wc.PageRequest(req)
		return req.WriteResponseHTML(data)
	})
	return nil
}
