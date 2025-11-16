package webclient

import (
	"embed"

	client "github.com/net12labs/cirm/dali/site-client-web"
)

//go:embed web/*
var content embed.FS

type ProviderClient struct {
	*client.Client
}

func NewWebClient() *ProviderClient {
	return &ProviderClient{Client: client.NewClient()}
}

func (wc *ProviderClient) Init() error {
	wc.Server.AddRoute("/provider", func(req *client.Request) error {
		data, err := content.ReadFile("web/index.html")
		if err != nil {
			return req.WriteResponse404()
		}
		return req.WriteResponseHTML(data)
	})
	return nil
}
