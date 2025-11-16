package service

import (
	"fmt"

	apiwebserver "github.com/net12labs/cirm/dali/api-web-server"
	clientwebserver "github.com/net12labs/cirm/dali/client-web-server"
	"github.com/net12labs/cirm/dali/context/service"
	"github.com/net12labs/cirm/dali/data"
	rtm "github.com/net12labs/cirm/dali/runtime"
	webserver "github.com/net12labs/cirm/dali/web-server"
	"github.com/net12labs/cirm/service-daemon/svc/admin"
	"github.com/net12labs/cirm/service-daemon/svc/platform"
	"github.com/net12labs/cirm/service-daemon/svc/provider"
	"github.com/net12labs/cirm/service-daemon/svc/root"
	"github.com/net12labs/cirm/service-daemon/svc/user"
)

type Serve struct {
	*service.Service
	User         *user.Unit
	Provider     *provider.Unit
	Platform     *platform.Unit
	Root         *root.Unit
	Admin        *admin.Unit
	ApiServer    *apiwebserver.Server
	ClientServer *clientwebserver.Server
}

// so the actual web server may need to be started even higher

func NewServe() *Serve {
	sv := &Serve{
		Service:  service.NewService(),
		User:     user.NewUnit(),
		Admin:    admin.NewUnit(),
		Provider: provider.NewUnit(),
		Root:     root.NewUnit(),
		Platform: platform.NewUnit(),
	}
	sv.WebServer = webserver.NewWebServer()
	sv.ApiServer = apiwebserver.NewServer()
	sv.ClientServer = clientwebserver.NewServer()

	sv.ApiServer.WebServer = sv.WebServer
	sv.ClientServer.WebServer = sv.WebServer

	sv.Root.WebClient.Server = sv.ClientServer
	sv.Provider.WebClient.Server = sv.ClientServer
	sv.Platform.WebClient.Server = sv.ClientServer
	sv.Admin.WebClient.Server = sv.ClientServer
	sv.User.WebClient.Server = sv.ClientServer

	sv.Root.WebApi.Server = sv.ApiServer
	sv.Provider.WebApi.Server = sv.ApiServer
	sv.Platform.WebApi.Server = sv.ApiServer
	sv.Admin.WebApi.Server = sv.ApiServer
	sv.User.WebApi.Server = sv.ApiServer

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
	dbPath := rtm.Etc.Get("main_db_path").String()
	rtm.Do.InitFsPath(dbPath)
	mainDb := data.Ops.CreateDb("main", dbPath)

	if err := mainDb.Init(); err != nil {
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
