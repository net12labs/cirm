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
	agt := &WebAiAgentApi{WebApi: webapi.NewWebApi()}
	agt.Domain.Path = "/platform/ai-agent/api"
	return agt
}

func (api *WebAiAgentApi) Init() {
	api.WebApi.Server.AddRoute(api.Domain.Path+"/refresh-data", func(req *webapi.Request) error {
		req.Response.StatusCode = http.StatusOK
		req.WriteResponse([]byte("Data refresh triggered"))
		return nil
	})
}
