package webclient

import (
	"embed"

	client "github.com/net12labs/cirm/dolly/site-client-web"
	domain "github.com/net12labs/cirm/domain-web"
)

//go:embed web/*
var content embed.FS

type WebClient struct {
	*client.Client
}

func NewWebClient() *WebClient {
	cl := &WebClient{Client: client.NewClient()}
	cl.Domain.Path = domain.UrlPrefix
	return cl
}

func HandleHtmlRequest(req *client.Request, filePath string) error {
	data, err := content.ReadFile(filePath)
	if err != nil {
		return req.WriteResponse404()
	}
	return req.WriteResponseHTML(data)
}

func (wc *WebClient) Init() error {
	wc.Server.AddRoute(wc.Domain.Path, func(req *client.Request) error {
		return HandleHtmlRequest(req, "web/guest-lander.html")
	})

	wc.Server.AddRoute(wc.Domain.MakePath("home"), func(req *client.Request) error {
		return HandleHtmlRequest(req, "web/home.html")
	})

	return nil
}
