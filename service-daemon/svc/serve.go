package service

import (
	"fmt"

	"github.com/net12labs/cirm/dali/context/service"
	rtm "github.com/net12labs/cirm/dali/runtime"
	"github.com/net12labs/cirm/service-daemon/data"
	"github.com/net12labs/cirm/service-daemon/svc/admin"
	"github.com/net12labs/cirm/service-daemon/svc/platform"
	"github.com/net12labs/cirm/service-daemon/svc/provider"
	"github.com/net12labs/cirm/service-daemon/svc/root"
	"github.com/net12labs/cirm/service-daemon/svc/user"
)

type Serve struct {
	*service.Service
	User     *user.Unit
	Provider *provider.Unit
	Platform *platform.Unit
	Root     *root.Unit
	Admin    *admin.Unit
}

func NewServe() *Serve {
	sv := &Serve{
		Service:  service.NewService(),
		User:     user.NewUnit(),
		Admin:    admin.NewUnit(),
		Provider: provider.NewUnit(),
		Root:     root.NewUnit(),
		Platform: platform.NewUnit(),
	}
	sv.WebServer = service.NewWebServer()
	sv.Root.WebServer = sv.WebServer
	sv.User.WebServer = sv.WebServer
	sv.Admin.WebServer = sv.WebServer
	sv.Provider.WebServer = sv.WebServer
	sv.Platform.WebServer = sv.WebServer

	return sv
}

func (s *Serve) Start() error {
	if err := s.runtimeInit(); err != nil {
		fmt.Println(err)
		rtm.Runtime.Exit(1)
	}

	if err := s.initClients(); err != nil {
		fmt.Println(err)
		rtm.Runtime.Exit(1)
	}
	rtm.Runtime.Exit(0)
	return nil
}

func (s *Serve) runtimeInit() error {

	pid := rtm.Pid
	pid.Pid.PidFilePath = rtm.Etc.Get("pid_file_path").String()
	if err := pid.Handle_ExitOnDuplicate(); err != nil {
		fmt.Println(err)
		rtm.Runtime.Exit(1)
	}
	if err := s.dataInit(); err != nil {
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

	return nil
}

func (s *Serve) dataInit() error {
	dbPath := rtm.Etc.Get("data_dir").String() + "/data.db"
	data.Module.DB.DbPath = dbPath
	rtm.Do.InitFsPath(data.Module.DB.DbPath)

	if err := data.Module.Init(); err != nil {
		fmt.Println("Failed to initialize database:", err)
		rtm.Runtime.Exit(1)
	}

	return nil
}

func (s *Serve) initClients() error {

	s.User.Mode.SetKeys("web", "cli")
	s.User.OnExit = func() {
		// Cleanup tasks here
	}
	if err := s.Root.Init(); err != nil {
		fmt.Println(err)
		rtm.Runtime.Exit(1)
	}

	if err := s.Admin.Init(); err != nil {
		fmt.Println(err)
		rtm.Runtime.Exit(1)
	}

	if err := s.Platform.Init(); err != nil {
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
