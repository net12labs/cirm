package domain

import (
	"fmt"

	"github.com/net12labs/cirm/ops/rtm"
)

type AstroDom struct {
	name     string
	Services *Svcs
	Admin    *domAdmin
	MainSite *WebSite
	StdOut   func(string)
	StdErr   func(string)
}

var Main = &AstroDom{Services: NewSvcs(), MainSite: NewWebSite()}

// we should also be able to bubble from here up

func (d *AstroDom) Init() *AstroDom {
	d.name = rtm.Etc.Get("domain_name").String()
	d.Services.Init()

	d.Admin = NewDomAdmin()

	d.Admin.WebServer.WebServer = d.Services.WebServer
	d.MainSite.Site.WebServer = d.Services.WebServer

	d.Admin.Init()
	d.MainSite.Init()

	return d
}

func (d *AstroDom) Start() {
	d.MainSite.Start()

	if err := d.Services.WebServer.Start(); err != nil {
		fmt.Printf("Failed to start web server: %v\n", err)
		rtm.Runtime.Exit(1)
	}
}
