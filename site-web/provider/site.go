package website

import (
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
