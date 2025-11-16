package root

import (
	"net/http"

	"github.com/net12labs/cirm/dali/context/webapi"
	"github.com/net12labs/cirm/dali/shell"
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

		rq, err := req.ReadRequestBodyAsMap()
		if err != nil {
			req.Response.StatusCode = http.StatusBadRequest
			req.WriteResponse(map[string]any{"message": "Invalid request"})
			return
		}

		username, _ := rq["username"].(string)
		password, _ := rq["password"].(string)

		if api.svc.Agent.UserLogin(username, password) == nil {
			req.Response.StatusCode = http.StatusOK
			req.WriteResponse(map[string]any{"token": "loonabalooona"})
		} else {
			req.Response.StatusCode = http.StatusUnauthorized
			req.WriteResponse(map[string]any{"message": "Login failed"})
		}
	})

	api.Server.AddRoute("/api/create-account", func(req *webapi.Request) {

		rq, err := req.ReadRequestBodyAsMap()
		if err != nil {
			req.Response.StatusCode = http.StatusBadRequest
			req.WriteResponse(map[string]any{"message": "Invalid request"})
			return
		}

		username := rq["username"].(string)
		password := rq["password"].(string)

		if username == "" || password == "" {
			req.Response.StatusCode = http.StatusBadRequest
			req.WriteResponse(map[string]any{"message": "Username and password are required"})
			return
		}

		shell := shell.NewShell(1) // assuming user ID 1 is root
		newAccount := shell.CreateStdAccount(username)
		shell.SetAccountPassword(newAccount.Id(), password)
		req.Response.StatusCode = http.StatusOK
		req.WriteResponse(map[string]any{"user_id": newAccount.Id, "username": newAccount.Name})
	})
}
