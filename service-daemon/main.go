package main

import (
	"fmt"

	"github.com/net12labs/cirm/dali/context/cmd"
	"github.com/net12labs/cirm/service-daemon/exec"
	"github.com/net12labs/cirm/service-daemon/meta"
	"github.com/net12labs/cirm/service-daemon/site"

	rtm "github.com/net12labs/cirm/dali/runtime"
)

// so context we can also package in the db - so it is all atomic
// this is the root level
// And then we track session  - as browser session and tab session
// so preferable we would be keeping a goroutine alive per session/sub session

func main() {
	rtm.Etc.SetKV("unit_id", "default")
	rtm.Etc.SetKV("rtm_name", "china-ip-routes-maker")
	rtm.Etc.SetKV("pid_file_path", "../units/"+rtm.Etc.Get("unit_id").String()+"/proc/china-ip-routes-maker.pid")
	rtm.Etc.SetKV("data_dir", "../units/"+rtm.Etc.Get("unit_id").String()+"/data")
	rtm.Etc.SetKV("main_db_path", rtm.Etc.Get("data_dir").String()+"/main.db")

	rtm.Runtime.OnPanic.AddListener(func(err any) {
		fmt.Println("Runtime Panic:", err)
	})

	if rtm.Args.HasKey("--run") {
		go func() {
			site := site.NewSite()

			site.Execute = func(cmd *cmd.Cmd) {
				fmt.Println("Executing command via Site:", cmd)
				if cmd.ExitCode == -1 {
					cmd.ExitCode = 1
					cmd.ErrorMsg = "No handler implemented"
				}
			}

			if err := site.Start(); err != nil {
				fmt.Println("Failed to start service:", err)
				rtm.Runtime.Exit(1)
			}
			rtm.Runtime.Exit(0)
		}()
		rtm.Runtime.WaitForSIGTERM()
	}

	if rtm.Args.HasKey("--exec") {
		cmd := exec.NewCmd()
		if err := cmd.Execute(); err != nil {
			fmt.Println("Failed to execute command:", err)
			rtm.Runtime.Exit(1)
		}
		rtm.Runtime.Exit(0)
	}

	if rtm.Args.HasKey("--help") || rtm.Args.HasKey("-h") {
		meta.PrintHelp()
		rtm.Runtime.Exit(0)
	}

	meta.PrintHelp()
	rtm.Runtime.Exit(0)

}
