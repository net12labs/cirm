package site

import "github.com/net12labs/cirm/dali/context/cmd"

type Site struct {
	Execute func(cmd *cmd.Cmd)
}

func NewSite() *Site {
	site := &Site{}
	// Initialize Site fields here
	return site
}

func (a *Site) OnExecute(cmd *cmd.Cmd) {
	a.Execute(cmd)
}
