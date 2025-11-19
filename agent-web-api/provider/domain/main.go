package domain

import (
	dom "github.com/net12labs/cirm/ops/domain"
)

func Domain() *Dom {
	return dom.Provider_Agent_Api
}

type Dom = dom.AgentApiDomain
