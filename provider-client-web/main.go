package providerclientweb

import (
	"embed"

	webserver "github.com/net12labs/cirm/dali/web-server"
)

//go:embed web/*
var content embed.FS

type ProviderClient struct {
	Server *webserver.WebServer
}

func NewProviderClient() *ProviderClient {
	return &ProviderClient{}
}

func (wc *ProviderClient) Init() error {
	wc.Server.AddRoute("/provider", func(req *webserver.Request) {
		// Serve the main page
		data, err := content.ReadFile("web/index.html")
		if err != nil {
			req.Response = &webserver.Response{
				StatusCode: 404,
			}
			req.WriteResponse([]byte("404 Not Found"))
			return
		}
		req.Response = &webserver.Response{
			StatusCode: 200,
			MimeType:   "text/html",
		}
		req.WriteResponse(data)
	})
	return nil
}
