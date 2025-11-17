package agent

import "github.com/net12labs/cirm/dali/context/cmd"

type Agent struct {
	Execute func(cmd *cmd.Cmd)
}

func NewAgent() *Agent {
	agent := &Agent{}
	// Initialize Agent fields here
	return agent
}
func (a *Agent) OnExecute(cmd *cmd.Cmd) {
	a.Execute(cmd)
}
