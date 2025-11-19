package domain

import (
	"fmt"

	site "github.com/net12labs/cirm/astro-host/host-web"
	"github.com/net12labs/cirm/dolly/cmd"
	"github.com/net12labs/cirm/ops/rtm"
)

type WebSite struct {
	Site *site.Site
}

func NewWebSite() *WebSite {
	return &WebSite{
		Site: site.NewSite(),
	}
}

type TaskShell struct {
}

func (d *WebSite) OnExecute(cmd *cmd.Cmd) {
	fmt.Println("Executing command via Site:", cmd)

	if cmd.Cmd == "domain.shutdown" {
		fmt.Println("Shutting down service...")
		rtm.Runtime.Exit(0)
		return
	}

	if cmd.Cmd == "domain.user.create" {
		// Example implementation for user creation
		username := cmd.Params["username"].(string)
		password := cmd.Params["password"].(string)
		fmt.Printf("Creating user: %s with password: %s\n", username, password)
		// Add actual user creation logic here
		cmd.ExitCode = 0
		return
	}

	d.Execute(cmd)
}

func (d *WebSite) Execute(cmd *cmd.Cmd) {
	// empty loopback set intentionally
}

func (d *WebSite) Start() {

	if err := d.Site.Start(); err != nil {
		fmt.Println("Failed to start service:", err)
		rtm.Runtime.Exit(1)
	}

}

func (d *WebSite) Init() {
	d.Site.Init()
}
