package main

import (
	"fmt"

	"github.com/net12labs/cirm/service-daemon/domain"
	"github.com/net12labs/cirm/service-daemon/exec"
	"github.com/net12labs/cirm/service-daemon/meta"

	"github.com/net12labs/cirm/dali/context/cmd"
	"github.com/net12labs/cirm/dali/rtm"
)

// so context we can also package in the db - so it is all atomic
// this is the root level
// And then we track session  - as browser session and tab session
// so preferable we would be keeping a goroutine alive per session/sub session
// Sessions can also be remote

func main() {

	rtm.Etc.SetKV("unit_id", "default")
	rtm.Etc.SetKV("rtm_name", "china-ip-routes-maker")
	rtm.Etc.SetKV("domain_name", rtm.Etc.GetJoined("/", "rtm_name", "unit_id"))
	rtm.Etc.SetKV("home_dir", "../units/"+rtm.Etc.Get("unit_id").String())
	rtm.Etc.SetKV("pid_file_path", rtm.Etc.Get("home_dir").String()+"/proc/china-ip-routes-maker.pid")
	rtm.Etc.SetKV("data_dir", rtm.Etc.Get("home_dir").String()+"/data")
	rtm.Etc.SetKV("main_db_path", rtm.Etc.Get("data_dir").String()+"/main.db")
	rtm.Etc.SetKV("socket_path", rtm.Etc.Get("home_dir").String()+"/main.sock")

	// so the domain needs to have separate logic (like a kernel)

	rtm.Runtime.OnPanic.AddListener(func(err any) {
		fmt.Println("Runtime Panic:", err)
	})

	if rtm.Args.HasKey("--run") {
		go func() {
			dom := domain.Main.Init(rtm.Etc.Get("domain_name").String())

			dom.Execute = func(cmd *cmd.Cmd) {
				if cmd.ExitCode == -1 {
					cmd.ExitCode = 1
					cmd.ErrorMsg = "No handler implemented"
				}
			}

			dom.Start()
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
