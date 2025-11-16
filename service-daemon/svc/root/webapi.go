package root

import (
	"net/http"

	"github.com/net12labs/cirm/dali/context/webapi"
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
	api.Server.AddRoute("/api/login", func(req *webapi.Request) {
		if api.svc.Agent.UserLogin("abc", "password") == nil {
			req.Response = &webapi.Response{
				StatusCode: http.StatusOK,
			}
			req.WriteResponse(map[string]any{"token": "loonabalooona"})
		} else {
			req.Response = &webapi.Response{
				StatusCode: http.StatusUnauthorized,
			}
			req.WriteResponse(map[string]any{"message": "Login failed"})
		}
	})

}
