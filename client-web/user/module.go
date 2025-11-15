package webclient

import (
	"embed"

	client "github.com/net12labs/cirm/dali/client-web"
	webserver "github.com/net12labs/cirm/dali/web-server"
)

//go:embed web/*
var content embed.FS

type WebClient struct {
	*client.Client
}

func NewWebClient() *WebClient {
	return &WebClient{}
}

func (wc *WebClient) OnRequest(req *webserver.Request) error {
	// Start the web client
	return nil
}

func (wc *WebClient) Init() error {
	wc.Server.AddRoute("/user", func(req *webserver.Request) {
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
