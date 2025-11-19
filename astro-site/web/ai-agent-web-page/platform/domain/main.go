package domain

import (
	dom "github.com/net12labs/cirm/ops/domain"
)

func Domain() *Dom {
	return dom.Platform_AiAgent_Page
}

type Dom = dom.AiAgentPageDomain
