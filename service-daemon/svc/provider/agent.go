package provider

import (
	"github.com/net12labs/cirm/dali/context/agent"
	refreshdata "github.com/net12labs/cirm/service-daemon/svc/provider/work/operation/refresh-data"
)

type SvcAgent struct {
	*agent.Agent
	Svc *Unit
}

func (sa *SvcAgent) RefreshData() {
	op := refreshdata.RefreshData{}
	op.Init()
	op.Op.Execute()
}
