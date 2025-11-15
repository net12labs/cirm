package webclient

import (
	"embed"

	"github.com/net12labs/cirm/dali/client-web"
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
	wc.Server.AddRoute("/provider", func(req *client.Request) {
		// Serve the main page
		data, err := content.ReadFile("web/index.html")
		if err != nil {
			req.Response = &client.Response{
				StatusCode: 404,
			}
			req.WriteResponse([]byte("404 Not Found"))
			return
		}
		req.Response = &client.Response{
			StatusCode: 200,
			MimeType:   "text/html",
		}
		req.WriteResponse(data)
	})
	return nil
}
