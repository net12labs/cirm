package webclient

import (
	"embed"

	webserver "github.com/net12labs/cirm/dali/web-server"
)

//go:embed web/*
var content embed.FS

type WebClient struct {
	Server *webserver.WebServer
}

func NewWebClient() *WebClient {
	return &WebClient{}
}

func (wc *WebClient) Init() error {
	wc.Server.AddRoute("/platform", func(req *webserver.Request) {
		// Serve the main page
		if req.Path.Path != "/platform" {
			req.Response = &webserver.Response{
				StatusCode: 404,
			}
			req.WriteResponse([]byte("404 Not Found"))
			return
		}
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
