package aiagentwebapi

import (
	"net/http"

	api "github.com/net12labs/cirm/dali/client-api/ai-agent"
)

type WebAiAgentApi struct {
	*api.ClientApi
	// WebApi fields here
}

func NewWebApi() *WebAiAgentApi {
	agt := &WebAiAgentApi{ClientApi: api.NewClient()}
	agt.Domain.Path = "/consumer/ai-agent/api"
	return agt
}

func (wi *WebAiAgentApi) Init() {
	wi.ClientApi.Server.AddRoute(wi.Domain.Path+"/refresh-data", func(req *api.Request) error {
		req.Response.StatusCode = http.StatusOK
		req.WriteResponse([]byte("Data refresh triggered"))
		return nil
	})
}
