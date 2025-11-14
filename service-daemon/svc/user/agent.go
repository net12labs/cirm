package user

import (
	refreshdata "github.com/net12labs/cirm/service-daemon/svc/user/operation/refresh-data"
)

type SvcAgent struct {
	Svc *UserUnit
	// SvcAgent fields here
}

func (sa *SvcAgent) RefreshData() {
	op := refreshdata.RefreshData{}
	op.Init()
	op.Op.Execute()
}
