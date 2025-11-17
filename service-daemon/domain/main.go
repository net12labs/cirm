package domain

import (
	"fmt"

	"github.com/net12labs/cirm/dali/rtm"
)

type dom struct {
	name     string
	Services *Svcs
	Admin    *domAdmin
	WebSite  *WebSite
	StdOut   func(string)
	StdErr   func(string)
}

var Main = &dom{}

// we should also be able to bubble from here up

func (d *dom) Init() *dom {
	d.name = rtm.Etc.Get("domain_name").String()
	d.Services = NewSvcs()
	d.Services.Init()

	d.Admin = NewDomAdmin()
	d.WebSite = NewWebSite()

	d.Admin.WebServer.WebServer = d.Services.WebServer
	d.WebSite.Site.WebServer = d.Services.WebServer
	d.Services.Ecdn.Server.WebServer = d.Services.WebServer

	if err := d.Services.Ecdn.Server.Init(); err != nil {
		fmt.Printf("Failed to initialize ECDN server: %v\n", err)
		rtm.Runtime.Exit(1)
	}
	d.Admin.Init()
	d.WebSite.Init()

	return d
}

func (d *dom) Start() {
	d.WebSite.Start()

	if err := d.Services.WebServer.Start(); err != nil {
		fmt.Printf("Failed to start web server: %v\n", err)
		rtm.Runtime.Exit(1)
	}
}
