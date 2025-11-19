package webapi

import (
	"fmt"
	"net/http"

	webapi "github.com/net12labs/cirm/dali/client-api/website"
	"github.com/net12labs/cirm/dali/context/cmd"
	"github.com/net12labs/cirm/dali/shell"
)

type WebApi struct {
	*webapi.ClientApi
	// WebApi fields here
}

func NewWebApi() *WebApi {
	api := &WebApi{ClientApi: webapi.NewClient()}
	api.Domain.Path = "/site/api"
	return api
}

func (api *WebApi) Init() {
	api.Server.AddRoute(api.Domain.MakePath("login"), func(req *webapi.Request) error {

		rq, err := req.ReadRequestBodyAsMap()
		if err != nil {
			return err
		}

		username, _ := rq["username"].(string)
		password, _ := rq["password"].(string)
		if username == "" || password == "" {
			return fmt.Errorf("username and password are required")
		}

		cmd := cmd.NewCmd()
		cmd.Target = "user.login"
		cmd.Params = map[string]any{"username": username, "password": password}

		cmd.ExitCode = 0
		cmd.Result = map[string]any{"message": "Login successful", "token": "loonabalooona"}

		if cmd.ExitCode != 0 {
			req.Response.StatusCode = http.StatusBadRequest
			return req.WriteResponse(map[string]any{"message": cmd.ErrorMsg})
		}

		// if api.svc.Agent.UserLogin(username, password) == nil {
		req.Response.StatusCode = http.StatusOK
		return req.WriteResponse(map[string]any{"token": "loonabalooona"})
		// } else {
		// 	req.Response.StatusCode = http.StatusUnauthorized
		// 	req.WriteResponse(map[string]any{"message": "Login failed"})
		// }
	})

	api.Server.AddRoute(api.Domain.MakePath("create-account"), func(req *webapi.Request) error {

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
