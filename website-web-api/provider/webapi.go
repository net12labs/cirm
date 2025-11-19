package webapi

import (
	"net/http"

	webapi "github.com/net12labs/cirm/dali/client-api/website"
	dom "github.com/net12labs/cirm/website-web-api/provider/domain"
)

type WebApi struct {
	*webapi.ClientApi
}

func NewWebApi() *WebApi {
	api := &WebApi{ClientApi: webapi.NewClient()}
	return api
}

func (api *WebApi) Init() {

	api.ClientApi.Server.AddRoute(dom.Domain().MakePath("refresh-data"), func(req *webapi.Request) error {
		req.Response.StatusCode = http.StatusOK
		req.WriteResponse([]byte("Data refresh triggered"))
		return nil
	})

	api.ClientApi.Server.AddRoute(dom.Domain().MakePath("get-routes"), func(req *webapi.Request) error {
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
