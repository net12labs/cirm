package provider

import (
	refreshdata "github.com/net12labs/cirm/service-daemon/svc/provider/work/operation/refresh-data"
)

type SvcAgent struct {
	Svc *Unit
	// SvcAgent fields here
}

func (sa *SvcAgent) RefreshData() {
	op := refreshdata.RefreshData{}
	op.Init()
	op.Op.Execute()
}
