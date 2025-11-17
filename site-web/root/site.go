package website

import (
	"github.com/net12labs/cirm/dali/context/cmd"
	"github.com/net12labs/cirm/dali/context/site"
)

type Site struct {
	*site.Site
}

func NewSite() *Site {
	ag := &Site{}
	ag.Site = site.NewSite()
	return ag
}
func (a *Site) Init() error {
	return nil
}

func (a *Site) OnExecute(cmd *cmd.Cmd) {
	if cmd.Target == "user.login" {
		cmd.ExitCode = 0
		cmd.Result = map[string]any{"message": "Login successful", "token": "loonabalooona"}
		return
	}
	a.Execute(cmd)
}
