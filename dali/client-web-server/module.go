package clientwebserver

import (
	"fmt"
	"net/http"

	webserver "github.com/net12labs/cirm/mali/web-server"
)

type Server struct {
	WebServer *webserver.WebServer
}

func NewServer() *Server {
	return &Server{}
}
func (s *Server) AddRoute(path string, handler func(req *Request) error) error {
	s.WebServer.AddRoute(path, func(req *webserver.Request) {
		apiReq := &Request{Request: req, Response: &Response{Response: req.Response}}

		if err := handler(apiReq); err != nil {
			apiReq.Response.StatusCode = http.StatusInternalServerError
			fmt.Println("Failed to handle API request:", err)
			apiReq.WriteResponse(map[string]any{"message": "Internal server error"})
			return
		}
	})
	return nil
}
