package webagentapi

import (
	"net/http"

	dom "github.com/net12labs/cirm/astro-site/agent-web-api/admin/domain"
	api "github.com/net12labs/cirm/dali/client-api/agent"
)

type WebAgentApi struct {
	*api.ClientApi
	// WebApi fields here
}

func NewWebApi() *WebAgentApi {
	agt := &WebAgentApi{ClientApi: api.NewClient()}
	return agt
}

func (wi *WebAgentApi) Init() {
	wi.Server.AddRoute(dom.Domain().MakePath("refresh-data"), func(req *api.Request) error {
		req.Response.StatusCode = http.StatusOK
		req.WriteResponse([]byte("Data refresh triggered"))
		return nil
	})

	wi.Server.AddRoute(dom.Domain().MakePath("get-routes"), func(req *api.Request) error {
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
