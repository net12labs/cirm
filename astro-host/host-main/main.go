package proc_main

import (
	"fmt"

	domain "github.com/net12labs/cirm/astro-dom/dom-web"
	hostadmin "github.com/net12labs/cirm/astro-host/host-admin"
	sitedb "github.com/net12labs/cirm/astro-host/host-main/db_schema/site"
	"github.com/net12labs/cirm/astro-host/host-main/ecdn"
	astrowebmain "github.com/net12labs/cirm/astro-host/host-web"
	astrowebadmin "github.com/net12labs/cirm/astro-host/host-web-admin"
	mdata "github.com/net12labs/cirm/mali/data"
	webserver "github.com/net12labs/cirm/mali/web-server"
	"github.com/net12labs/cirm/ops/data"
	"github.com/net12labs/cirm/ops/rtm"
)

func processInit() {

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

type Main struct {
	WebServer    *webserver.WebServer
	HostWeb      *astrowebmain.Site
	HostWebAdmin *astrowebadmin.WebAdmin
	HostAdmin    *hostadmin.HostAdmin
	Domain       *domain.AstroDom
	Ecdn         *ecdn.Ecdn
}

func (m *Main) Run() {
	m.HostAdmin = hostadmin.NewHostAdmin()

	m.WebServer = webserver.NewWebServer()
	dbInit()
	m.Ecdn = ecdn.NewEcdn()
	m.WebServer.AddRoute("/", func(req *webserver.Request) {
		if req.Path.Path == "/" {
			req.RedirectToUrl("/site/site/home")
			return
		}
		req.WriteResponse404()
	})

	m.Ecdn.Server.WebServer = m.WebServer

	if err := m.Ecdn.Server.Init(); err != nil {
		fmt.Printf("Failed to initialize ECDN server: %v\n", err)
		rtm.Runtime.Exit(1)
	}

	m.Domain = domain.Main
	m.Domain.Services.WebServer = m.WebServer

	m.HostWeb = astrowebmain.NewSite()
	m.Domain.MainSite.Site = m.HostWeb
	m.HostWebAdmin = astrowebadmin.NewWebAdmin()

	dom := m.Domain.Init()
	dom.StdOut = func(msg string) {
		fmt.Println("STDOUT:", msg)
	}
	dom.StdErr = func(msg string) {
		fmt.Println("STDERR:", msg)
	}
	dom.Start()
}

func dbInit() {
	mainDb := data.Ops.CreateDb("dom", rtm.Etc.Get("dom_db_path").String())
	hostDb := data.Ops.CreateDb("host", rtm.Etc.Get("host_db_path").String())
	siteDb := data.Ops.CreateDb("site", rtm.Etc.Get("site_db_path").String())

	databases := []*mdata.SqliteDb{mainDb, hostDb, siteDb}
	for _, db := range databases {
		if err := db.Init(); err != nil {
			fmt.Println("Failed to initialize database:", err)
			rtm.Runtime.Exit(1)
		}
	}
	sitedb.InitSchema()
	sitedb.Dbo.Db = siteDb

	fmt.Println("=== Resetting database ===")
	if err := sitedb.Dbo.InitDb(); err != nil {
		fmt.Println("Failed to reset database:", err)
		rtm.Runtime.Exit(1)
	}
	fmt.Println("=== Database initialized successfully ===")
	fmt.Printf("Tables created: %v\n", sitedb.Dbo.ListTables())

}

func Try() {

	rtm.Runtime.OnPanic.AddListener(func(err any) {
		fmt.Println("Runtime Panic:", err)
	})

	if rtm.Args.HasKey("--run") {
		processInit()

		go func() {
			m := &Main{}
			m.Run()
			rtm.Runtime.Exit(0)
		}()
		rtm.Runtime.WaitForSIGTERM()
	}

}
