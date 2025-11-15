package provider

import (
	"github.com/net12labs/cirm/dali/context/cliapi"
)

type CliApi struct {
	*cliapi.CliApi
	svc *Unit

	// CliApi fields here
}

func (api *CliApi) Init() {

}
