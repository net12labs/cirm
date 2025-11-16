package site

import (
	"fmt"

	"github.com/net12labs/cirm/dali/context/cmd"
	"github.com/net12labs/cirm/dali/data"
	domain_context "github.com/net12labs/cirm/dali/domain/context"
	"github.com/net12labs/cirm/dali/rtm"
	apiwebserver "github.com/net12labs/cirm/dali/web-server-api"
	clientwebserver "github.com/net12labs/cirm/dali/web-server-client"
	"github.com/net12labs/cirm/service-daemon/site/admin"
	"github.com/net12labs/cirm/service-daemon/site/platform"
	"github.com/net12labs/cirm/service-daemon/site/provider"
	"github.com/net12labs/cirm/service-daemon/site/root"
	"github.com/net12labs/cirm/service-daemon/site/user"
)

type Site struct {
	*domain_context.Domain
	User         *user.Unit
	Provider     *provider.Unit
	Platform     *platform.Unit
	Root         *root.Unit
	Admin        *admin.Unit
	ApiServer    *apiwebserver.Server
	ClientServer *clientwebserver.Server
}

// so the actual web server may need to be started even higher

func NewSite() *Site {
	sv := &Site{
		Domain:   domain_context.NewDomain(),
		User:     user.NewUnit(),
		Admin:    admin.NewUnit(),
		Provider: provider.NewUnit(),
		Root:     root.NewUnit(),
		Platform: platform.NewUnit(),
	}
	sv.WebServer = domain_context.NewServer()
	sv.ApiServer = apiwebserver.NewServer()
	sv.ClientServer = clientwebserver.NewServer()

	sv.ApiServer.WebServer = sv.WebServer
	sv.ClientServer.WebServer = sv.WebServer

	sv.Root.WebClient.Server = sv.ClientServer
	sv.Provider.WebClient.Server = sv.ClientServer
	sv.Platform.WebClient.Server = sv.ClientServer
	sv.Admin.WebClient.Server = sv.ClientServer
	sv.User.WebClient.Server = sv.ClientServer

	sv.Root.WebAgentClient.Server = sv.ClientServer
	sv.Platform.WebAgentClient.Server = sv.ClientServer
	sv.Provider.WebAgentClient.Server = sv.ClientServer
	sv.Admin.WebAgentClient.Server = sv.ClientServer
	sv.User.WebAgentClient.Server = sv.ClientServer

	sv.Root.WebApi.Server = sv.ApiServer
	sv.Root.Execute = func(cmd *cmd.Cmd) {
		fmt.Println("Executing command via WebApi:", cmd)

		// User login should be handled right at the API edge
		// user create - that needs to be done by root user

		if cmd.Target == "user.login" {
			cmd.ExitCode = 0
			cmd.Result = map[string]any{"message": "Login successful", "token": "loonabalooona"}
			return
		}
		sv.Execute(cmd)
	}

	sv.Provider.WebApi.Server = sv.ApiServer
	sv.Platform.WebApi.Server = sv.ApiServer
	sv.Admin.WebApi.Server = sv.ApiServer
	sv.User.WebApi.Server = sv.ApiServer

	sv.Root.WebAgentApi.Server = sv.ApiServer
	sv.Platform.WebAgentApi.Server = sv.ApiServer
	sv.Provider.WebAgentApi.Server = sv.ApiServer
	sv.Admin.WebAgentApi.Server = sv.ApiServer
	sv.User.WebAgentApi.Server = sv.ApiServer

	return sv
}

func (s *Site) Start() error {

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

func (s *Site) runtimeInit() error {

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

func (s *Site) dataInit() error {
	dbPath := rtm.Etc.Get("main_db_path").String()
	rtm.Do.InitFsPath(dbPath)
	mainDb := data.Ops.CreateDb("main", dbPath)

	if err := mainDb.Init(); err != nil {
		fmt.Println("Failed to initialize database:", err)
		rtm.Runtime.Exit(1)
	}

	return nil
}

func (s *Site) initClients() error {

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
