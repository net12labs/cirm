package root

import (
	"fmt"
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
	api.Server.AddRoute("/api/login", func(req *webapi.Request) error {

		rq, err := req.ReadRequestBodyAsMap()
		if err != nil {
			return err
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
		return nil
	})

	api.Server.AddRoute("/api/create-account", func(req *webapi.Request) error {

		rq, err := req.ReadRequestBodyAsMap()
		if err != nil {
			return err
		}

		username := rq["username"].(string)
		password := rq["password"].(string)

		if username == "" || password == "" {
			return fmt.Errorf("username and password are required")
		}
		println("Creating account for:", username)

		shell := shell.NewShell(1) // assuming user ID 1 is root
		newAccount := shell.CreateStdAccount(username)
		shell.SetAccountPassword(newAccount.Id(), password)
		req.Response.StatusCode = http.StatusOK
		if err := req.WriteResponse(map[string]any{"user_id": newAccount.Id(), "username": newAccount.Name()}); err != nil {
			fmt.Println("Failed to write response:", err)
			return err
		}
		return nil
	})
}
