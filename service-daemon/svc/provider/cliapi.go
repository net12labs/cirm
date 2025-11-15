package provider

import (
	cliserver "github.com/net12labs/cirm/dali/cli-server"
)

type CliApi struct {
	svc    *Unit
	Server *cliserver.CliServer

	// CliApi fields here
}

func (api *CliApi) Init() {
	api.Server = &cliserver.CliServer{}
	api.Server.Init()
	api.Server.AddRoute("/refresh-data", func() {
		api.svc.Agent.RefreshData()
	})
	api.Server.AddRoute("/get-routes", func() {
		// we can have format like bash script, bird config, etc.
	})
}
