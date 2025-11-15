package admin

import (
	"net/http"

	webserver "github.com/net12labs/cirm/dali/web-server"
)

type WebApi struct {
	svc    *Unit
	Server *webserver.WebServer
	// WebApi fields here
}

func NewWebApi() *WebApi {
	return &WebApi{}
}

func (api *WebApi) Init() {
	api.Server.AddRoute("/admin/api/refresh-data", func(req *webserver.Request) {
		req.Response = &webserver.Response{
			StatusCode: http.StatusOK,
		}
		req.WriteResponse([]byte("Data refresh triggered"))
		api.svc.Agent.RefreshData()
	})

	api.Server.AddRoute("/admin/api/get-routes", func(req *webserver.Request) {
		// Get format from query parameter (bash, bird, json, etc.)
		format := req.Req.URL.Query().Get("format")
		if format == "" {
			format = "json"
		}

		// Example response structure
		response := map[string]interface{}{
			"format": format,
			"routes": []string{
				"route 1.1.1.0/24 via 192.168.1.1",
				"route 8.8.8.0/24 via 192.168.1.1",
			},
			"count": 2,
		}

		req.Response = &webserver.Response{
			StatusCode: http.StatusOK,
			MimeType:   "application/json",
		}
		req.WriteResponse(response)
	})
}
