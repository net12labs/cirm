package service

import (
	"fmt"

	"github.com/net12labs/cirm/service-daemon/data"

	ox "github.com/net12labs/cirm/dali/ox"
	rtm "github.com/net12labs/cirm/dali/runtime"
	webserver "github.com/net12labs/cirm/dali/web-server"

	webclient "github.com/net12labs/cirm/client-web/root"
	admin "github.com/net12labs/cirm/service-daemon/svc/admin"
	"github.com/net12labs/cirm/service-daemon/svc/provider"
	user "github.com/net12labs/cirm/service-daemon/svc/user"
)

type Serve struct {
	// Serve fields here
	WebServer *webserver.WebServer
	User      *user.Unit
	Provider  *provider.Unit
	Admin     *admin.Unit
	WebClient *webclient.WebClient
}

func NewServe() *Serve {
	sv := &Serve{
		User:      user.NewUnit(),
		Admin:     admin.NewUnit(),
		Provider:  provider.NewUnit(),
		WebServer: webserver.NewWebServer(),
		WebClient: webclient.NewWebClient(),
	}
	sv.WebClient.Server = sv.WebServer
	sv.User.WebServer = sv.WebServer
	sv.Admin.WebServer = sv.WebServer
	sv.Provider.WebServer = sv.WebServer

	return sv
}

func (s *Serve) Start() error {

	pid := ox.NewPidHandler().Init()
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

	data.Module.DB.DbPath = ox.Etc.Get("data_dir").String() + "/" + ox.Etc.Get("unit_id").String() + "/data/data.db"
	rtm.Do.InitFsPath(data.Module.DB.DbPath)

	if err := data.Module.Init(); err != nil {
		fmt.Println("Failed to initialize database:", err)
		rtm.Runtime.Exit(1)
	}

	if err := s.initClients(); err != nil {
		fmt.Println(err)
		rtm.Runtime.Exit(1)
	}
	rtm.Runtime.Exit(0)
	return nil

}

func (s *Serve) initClients() error {

	s.User.Mode.SetKeys("web", "cli")
	s.User.OnExit = func() {
		// Cleanup tasks here
	}
	if err := s.WebClient.Init(); err != nil {
		fmt.Println(err)
		rtm.Runtime.Exit(1)
	}

	if err := s.Admin.Init(); err != nil {
		fmt.Println(err)
		rtm.Runtime.Exit(1)
	}

	if err := s.Provider.Init(); err != nil {
		fmt.Println(err)
		rtm.Runtime.Exit(1)
	}

	if err := s.User.Init(); err != nil {
		fmt.Println(err)
		rtm.Runtime.Exit(1)
	}

	if err := s.WebServer.Start(); err != nil {
		fmt.Printf("Failed to start web server: %v\n", err)
		return err
	}
	return nil

}
