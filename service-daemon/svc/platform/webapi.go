package platform

import (
	"net/http"

	webapi "github.com/net12labs/cirm/dali/context/webapi"
)

type WebApi struct {
	*webapi.WebApi
	svc *Unit
	// WebApi fields here
}

func NewWebApi() *WebApi {
	return &WebApi{WebApi: webapi.NewWebApi()}
}

func (api *WebApi) Init() {
	api.WebApi.Server.AddRoute("/platform/api/refresh-data", func(req *webapi.Request) {
		req.Response.StatusCode = http.StatusOK
		req.WriteResponse([]byte("Data refresh triggered"))
		api.svc.Agent.RefreshData()
	})

	api.WebApi.Server.AddRoute("/platform/api/get-routes", func(req *webapi.Request) {
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

		req.Response.StatusCode = http.StatusOK
		req.Response.MimeType = "application/json"
		req.WriteResponse(response)
	})
}
