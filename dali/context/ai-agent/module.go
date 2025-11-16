package aiagent

import "github.com/net12labs/cirm/dali/context/cmd"

type AiAgent struct {
	Execute   func(cmd *cmd.Cmd)
	OnExecute func(cmd *cmd.Cmd)
}

func NewAiAgent() *AiAgent {
	agent := &AiAgent{}
	// Initialize Agent fields here
	return agent
}
