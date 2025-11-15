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
	wc.Server.AddRoute("/", func(req *client.Request) {
		// Serve the main page
		if req.Path.Path != "/" {
			req.Response = &client.Response{
				StatusCode: 404,
			}
			req.WriteResponse([]byte("404 Not Found"))
			return
		}
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
