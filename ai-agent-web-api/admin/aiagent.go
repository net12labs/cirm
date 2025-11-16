package aiagentwebapi

import (
	"net/http"

	webapi "github.com/net12labs/cirm/dali/context/webapi"
)

type WebAiAgentApi struct {
	*webapi.WebApi
	// WebApi fields here
}

func NewWebApi() *WebAiAgentApi {
	return &WebAiAgentApi{WebApi: webapi.NewWebApi()}
}

func (api *WebAiAgentApi) Init() {
	api.WebApi.Server.AddRoute("/admin/ai-agent/api/refresh-data", func(req *webapi.Request) error {
		req.Response.StatusCode = http.StatusOK
		req.WriteResponse([]byte("Data refresh triggered"))
		return nil
	})
}
