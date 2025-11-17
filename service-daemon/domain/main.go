package domain

import (
	"fmt"

	"github.com/net12labs/cirm/dali/context/cmd"
	"github.com/net12labs/cirm/dali/data"
	domain_context "github.com/net12labs/cirm/dali/domain/context"
	"github.com/net12labs/cirm/dali/rtm"
	socketserver "github.com/net12labs/cirm/dali/socket-server"
	webserver "github.com/net12labs/cirm/mali/web-server"
	"github.com/net12labs/cirm/service-daemon/site"
)

type dom struct {
	name         string
	Site         *site.Site
	WebServer    *domain_context.WebServer
	SocketServer *socketserver.SocketServer
	Execute      func(cmd *cmd.Cmd)
}

var Main = &dom{}

func (d *dom) Name() string {
	return d.name
}

// we should also be able to bubble from here up

func (d *dom) OnExecute(cmd *cmd.Cmd) {
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

func (d *dom) Init(name string) *dom {
	d.name = name
	d.runtimeInit()

	d.dataInit()
	d.Site = site.NewSite()

	d.Site.Execute = func(cmd *cmd.Cmd) {
		d.OnExecute(cmd)
	}

	d.SocketServer = socketserver.NewSocketServer()
	socketPath := rtm.Etc.Get("socket_path").String()
	if err := d.SocketServer.Start(socketPath); err != nil {
		fmt.Println("Failed to start socket server:", err)
		rtm.Runtime.Exit(1)
	}

	d.WebServer = webserver.NewWebServer()
	d.Site.WebServer = d.WebServer

	d.WebServer.AddRoute("/", func(req *webserver.Request) {
		if req.Path.Path == "/" {
			req.RedirectToUrl("/site")
			return
		}
		req.WriteResponse404()
	})

	d.Site.Init()

	return d
}

func (d *dom) Start() {

	if err := d.Site.Start(); err != nil {
		fmt.Println("Failed to start service:", err)
		rtm.Runtime.Exit(1)
	}

	if err := d.WebServer.Start(); err != nil {
		fmt.Printf("Failed to start web server: %v\n", err)
		rtm.Runtime.Exit(1)
	}

}

func (d *dom) dataInit() {
	dbPath := rtm.Etc.Get("main_db_path").String()
	rtm.Do.InitFsPath(dbPath)
	mainDb := data.Ops.CreateDb("main", dbPath)

	if err := mainDb.Init(); err != nil {
		fmt.Println("Failed to initialize database:", err)
		rtm.Runtime.Exit(1)
	}
}

func (d *dom) runtimeInit() {

	pid := rtm.Pid
	pid.Pid.PidFilePath = rtm.Etc.Get("pid_file_path").String()
	if err := pid.Handle_ExitOnDuplicate(); err != nil {
		fmt.Println(err)
		rtm.Runtime.Exit(1)
	}
	rtm.Runtime.OnPanic.AddListener(func(err any) {
		pid.Handle_CleanupOnExit()
	})
	rtm.Runtime.OnExit.AddListener(func(code any) {
		pid.Handle_CleanupOnExit()
		fmt.Println("Exited with code", code)
	})

}
