package webclient

import (
	"embed"

	"github.com/net12labs/cirm/dali/client-web"
)

//go:embed web/*
var content embed.FS

type WebClient struct {
	*client.Client
}

func NewWebClient() *WebClient {
	return &WebClient{Client: client.NewClient()}
}

func (wc *WebClient) Init() error {
	wc.Server.AddRoute("/", func(req *client.Request) error {
		if req.Path.Path != "/" {
			return req.WriteResponse404()
		}
		data, err := content.ReadFile("web/index.html")
		if err != nil {
			return req.WriteResponse404()
		}
		return req.WriteResponseHTML(data)
	})

	return nil
}
