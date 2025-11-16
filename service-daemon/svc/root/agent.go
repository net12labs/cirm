package root

import (
	"github.com/net12labs/cirm/dali/context/agent"
)

type SvcAgent struct {
	*agent.Agent
	Svc *Unit
	// SvcAgent fields here
}

func (sa *SvcAgent) UserLogin(username, password string) error {
	// Implement user login logic here
	return nil
}
