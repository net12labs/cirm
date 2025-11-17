package webagent

import (
	"github.com/net12labs/cirm/dali/context/agent"
)

type Agent struct {
	*agent.Agent
}

func (a *Agent) RefreshData() {

}

func NewAgent() *Agent {
	ag := &Agent{}
	ag.Agent = agent.NewAgent()
	return ag
}

func (a *Agent) Init() error {
	return nil
}
