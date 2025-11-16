package aiagentweb

import (
	aiagent "github.com/net12labs/cirm/dali/context/ai-agent"
)

type AiAgent struct {
	*aiagent.AiAgent
}

func (a *AiAgent) RefreshData() {

}

func NewAiAgent() *AiAgent {
	ag := &AiAgent{}
	ag.AiAgent = aiagent.NewAiAgent()
	return ag
}
func (a *AiAgent) Init() error {
	return nil
}
