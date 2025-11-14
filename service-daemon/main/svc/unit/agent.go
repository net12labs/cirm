package unit

import (
	"cirm/svc/operation/refresh-data"
)

type SvcAgent struct {
	Svc *SvcUnit
	// SvcAgent fields here
}

func (sa *SvcAgent) RefreshData() {
	op := refreshdata.RefreshData{}
	op.Init()
	op.Op.Execute()
}
