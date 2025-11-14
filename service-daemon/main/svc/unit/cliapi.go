package unit

import (
	cliserver "github.com/net12labs/cirm/service-daemon/mod/cli-server"
)

type CliApi struct {
	svc    *SvcUnit
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
