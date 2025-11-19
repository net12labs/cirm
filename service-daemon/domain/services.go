package domain

import (
	"fmt"

	domain_context "github.com/net12labs/cirm/dali/domain/context"
	socketserver "github.com/net12labs/cirm/dolly/socket-server"
	webserver "github.com/net12labs/cirm/mali/web-server"
	"github.com/net12labs/cirm/ops/rtm"
	"github.com/net12labs/cirm/service-daemon/domain/ecdn"
)

type Svcs struct {
	WebServer    *domain_context.WebServer
	SocketServer *socketserver.SocketServer
	Ecdn         *ecdn.Ecdn
}

func NewSvcs() *Svcs {
	return &Svcs{}
}

func (s *Svcs) Init() {
	s.WebServer = webserver.NewWebServer()
	s.SocketServer = socketserver.NewSocketServer()
	s.Ecdn = ecdn.NewEcdn()

	socketPath := rtm.Etc.Get("socket_path").String()
	if err := s.SocketServer.Start(socketPath); err != nil {
		fmt.Println("Failed to start socket server:", err)
		rtm.Runtime.Exit(1)
	}

	s.WebServer.AddRoute("/", func(req *webserver.Request) {
		if req.Path.Path == "/" {
			req.RedirectToUrl("/site/home")
			return
		}
		req.WriteResponse404()
	})

}
