package web_server

import (
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
		err := handler(apiReq)
		if err != nil {
			req.Response.StatusCode = http.StatusInternalServerError
			req.WriteResponse([]byte("Internal Server Error"))
		}
	})
	return nil
}
