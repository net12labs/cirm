package user

import (
	"net/http"

	"github.com/net12labs/cirm/dali/context/webapi"
)

type WebApi struct {
	*webapi.WebApi
	svc *Unit
	// WebApi fields here
}

func (api *WebApi) Init() {

	api.Server.AddRoute("/user/api/refresh-data", func(req *webapi.Request) error {
		req.Response.StatusCode = http.StatusOK
		req.WriteResponse([]byte("Data refresh triggered"))
		api.svc.Agent.RefreshData()
		return nil
	})

	api.Server.AddRoute("/user/api/get-routes", func(req *webapi.Request) error {
		// Get format from query parameter (bash, bird, json, etc.)
		format := req.Req.URL.Query().Get("format")
		if format == "" {
			format = "json"
		}

		// Example response structure
		response := map[string]any{
			"format": format,
			"routes": []string{
				"route 1.1.1.0/24 via 192.168.1.1",
				"route 8.8.8.0/24 via 192.168.1.1",
			},
			"count": 2,
		}

		req.Response.StatusCode = http.StatusOK
		req.Response.MimeType = "application/json"
		return req.WriteResponse(response)
	})

}

func (api *WebApi) Start() {
}

func NewWebApi() *WebApi {
	return &WebApi{WebApi: webapi.NewWebApi()}
}
