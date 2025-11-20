package webapi

import (
	"fmt"
	"net/http"

	args_auth "github.com/net12labs/cirm/astro-site/args/auth"
	"github.com/net12labs/cirm/astro-site/svc"
	dom "github.com/net12labs/cirm/astro-site/web/website-web-api/site/domain"
	webapi "github.com/net12labs/cirm/dali/client-api/website"
)

type WebApi struct {
	*webapi.ClientApi
	// WebApi fields here
}

func NewWebApi() *WebApi {
	api := &WebApi{ClientApi: webapi.NewClient()}
	return api
}

func (api *WebApi) Init() {
	api.Server.AddRoute(dom.Domain().MakePath("login"), func(req *webapi.Request) error {

		creds := &args_auth.LoginCredentials{}
		if err := req.Json2Struct(creds); err != nil {
			return err
		}
		if !creds.IsValid() {
			return fmt.Errorf("username and password are required")
		}

		token, err := svc.Auth.UserLogin(creds)
		if err != nil {
			req.Response.StatusCode = http.StatusUnauthorized
			return req.WriteResponse(map[string]any{"message": "Login failed"})
		}
		req.Response.StatusCode = http.StatusOK
		return req.WriteResponse(map[string]any{"token": token, "redirectUrl": "/platform/site/home"})
	})

	api.Server.AddRoute(dom.Domain().MakePath("create-account"), func(req *webapi.Request) error {

		creds := &args_auth.LoginCredentials{}
		if err := req.Json2Struct(creds); err != nil {
			return err
		}
		if !creds.IsValid() {
			return fmt.Errorf("username and password are required")
		}

		if err := svc.Auth.AccountCreate(creds); err != nil {
			return fmt.Errorf("failed to create account: %w", err)
		}
		req.Response.StatusCode = http.StatusOK
		if err := req.WriteResponse(map[string]string{"status": "success"}); err != nil {
			fmt.Println("Failed to write response:", err)
			return err
		}

		return nil
	})
}
