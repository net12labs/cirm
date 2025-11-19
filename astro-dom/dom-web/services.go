package domain

import (
	"fmt"

	domain_context "github.com/net12labs/cirm/dali/domain/context"
	socketserver "github.com/net12labs/cirm/dolly/socket-server"
	"github.com/net12labs/cirm/ops/rtm"
)

type Svcs struct {
	WebServer    *domain_context.WebServer
	SocketServer *socketserver.SocketServer
}

func NewSvcs() *Svcs {
	return &Svcs{}
}

func (s *Svcs) Init() {
	s.SocketServer = socketserver.NewSocketServer()

	socketPath := rtm.Etc.Get("socket_path").String()
	if err := s.SocketServer.Start(socketPath); err != nil {
		fmt.Println("Failed to start socket server:", err)
		rtm.Runtime.Exit(1)
	}

}
