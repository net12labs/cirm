package platform

import (
	"github.com/net12labs/cirm/dali/context/agent"
)

type SvcAgent struct {
	*agent.Agent
	Svc *Unit
	// SvcAgent fields here
}

func (sa *SvcAgent) RefreshData() {

}
