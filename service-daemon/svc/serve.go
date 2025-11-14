package service

import (
	"fmt"

	"github.com/net12labs/cirm/service-daemon/data"

	ox "github.com/net12labs/cirm/dali/ox"
	rtm "github.com/net12labs/cirm/dali/runtime"
	webserver "github.com/net12labs/cirm/dali/web-server"

	admin "github.com/net12labs/cirm/service-daemon/svc/admin"
	user "github.com/net12labs/cirm/service-daemon/svc/user"
)

type Serve struct {
	// Serve fields here
	Server *webserver.WebServer
	User   *user.UserUnit
	Admin  *admin.AdminUnit
}

func NewServe() *Serve {
	sv := &Serve{
		User:  user.NewUserUnit(),
		Admin: admin.NewAdminUnit(),
	}
	sv.User.WebServer = sv.Server
	sv.Admin.WebServer = sv.Server
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

	s.User.Mode.SetKeys("web", "cli")
	s.User.OnExit = func() {
		// Cleanup tasks here
	}

	if err := s.User.Init(); err != nil {
		fmt.Println(err)
		rtm.Runtime.Exit(1)
	}

	exit_code := s.User.Run()

	if exit_code != 0 {
		fmt.Println(s.User.ExitMessage)
		rtm.Runtime.Exit(1)
	}
	rtm.Runtime.Exit(0)
	return nil

}
