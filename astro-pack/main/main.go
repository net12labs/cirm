package proc_main

import (
	"fmt"

	domain "github.com/net12labs/cirm/astro-dom/web-main"
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

func Try() {

	rtm.Runtime.OnPanic.AddListener(func(err any) {
		fmt.Println("Runtime Panic:", err)
	})

	if rtm.Args.HasKey("--run") {
		processInit()

		go func() {
			dom := domain.Main.Init()
			dom.StdOut = func(msg string) {
				fmt.Println("STDOUT:", msg)
			}
			dom.StdErr = func(msg string) {
				fmt.Println("STDERR:", msg)
			}
			dom.Start()
			rtm.Runtime.Exit(0)
		}()
		rtm.Runtime.WaitForSIGTERM()
	}

}
