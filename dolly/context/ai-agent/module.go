package aiagent

import "github.com/net12labs/cirm/dolly/context/cmd"

type AiAgent struct {
	Execute func(cmd *cmd.Cmd)
}

func NewAiAgent() *AiAgent {
	agent := &AiAgent{}
	// Initialize Agent fields here
	return agent
}

func (a *AiAgent) OnExecute(cmd *cmd.Cmd) {
	a.Execute(cmd)
}
