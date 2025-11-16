package apiwebserver

import (
	"net/http"

	webserver "github.com/net12labs/cirm/dali/web-server"
)

type Server struct {
	Server *webserver.WebServer
}

func NewServer() *Server {
	return &Server{}
}
func (s *Server) AddRoute(path string, handler func(req *ApiRequest)) error {
	s.Server.AddRoute(path, func(req *webserver.Request) {
		apiReq := &ApiRequest{Request: req, Response: &ApiResponse{Response: req.Response}}
		if !apiReq.Validate_HasBody() {
			apiReq.Response.StatusCode = http.StatusBadRequest
			apiReq.WriteResponse(map[string]any{"message": "Request body is required"})
			return
		}
		handler(apiReq)
	})
	return nil
}
