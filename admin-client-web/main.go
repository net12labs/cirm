package adminclientweb

import (
	"embed"

	webserver "github.com/net12labs/cirm/dali/web-server"
)

//go:embed web/*
var content embed.FS

type AdminClient struct {
	Server *webserver.WebServer
}

func NewAdminClient() *AdminClient {
	return &AdminClient{}
}

func (wc *AdminClient) Init() error {
	wc.Server.AddRoute("/admin", func(req *webserver.Request) {
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
