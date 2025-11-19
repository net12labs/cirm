package aiagentwebapi

import (
	"net/http"

	dom "github.com/net12labs/cirm/astro-site/ai-agent-web-api/provider/domain"
	api "github.com/net12labs/cirm/dali/client-api/ai-agent"
)

type WebAiAgentApi struct {
	*api.ClientApi
	// WebApi fields here
}

func NewWebApi() *WebAiAgentApi {
	agt := &WebAiAgentApi{ClientApi: api.NewClient()}
	return agt
}

func (wi *WebAiAgentApi) Init() {
	wi.ClientApi.Server.AddRoute(dom.Domain().Path(), func(req *api.Request) error {
		req.Response.StatusCode = http.StatusOK
		req.WriteResponse([]byte("Data refresh triggered"))
		return nil
	})
}
