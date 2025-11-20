package svc

import (
	"github.com/net12labs/cirm/astro-site/svc/auth"
)

var Auth *auth.Auth

func StartServices() {
	Auth = auth.NewService()
}
