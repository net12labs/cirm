package xmain

import (
	"fmt"

	"github.com/net12labs/cirm/service-daemon/data"

	ox "github.com/net12labs/cirm/dali/ox"
	rtm "github.com/net12labs/cirm/dali/runtime"

	svc "github.com/net12labs/cirm/service-daemon/main/svc"
)

type Serve struct {
	// Serve fields here
}

func NewServe() *Serve {
	return &Serve{}
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

	svc.Svc.Mode.SetKeys("web", "cli")
	svc.Svc.OnExit = func() {
		// Cleanup tasks here
	}

	if err := svc.Svc.Init(); err != nil {
		fmt.Println(err)
		rtm.Runtime.Exit(1)
	}

	exit_code := svc.Svc.Run()

	if exit_code != 0 {
		fmt.Println(svc.Svc.ExitMessage)
		rtm.Runtime.Exit(1)
	}
	rtm.Runtime.Exit(0)
	return nil

}
