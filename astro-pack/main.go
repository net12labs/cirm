package main

import (
	"fmt"

	"github.com/net12labs/cirm/astro-pack/domain"
	"github.com/net12labs/cirm/astro-pack/exec"
	"github.com/net12labs/cirm/astro-pack/meta"

	"github.com/net12labs/cirm/ops/rtm"
)

// so context we can also package in the db - so it is all atomic
// this is the root level
// And then we track session  - as browser session and tab session
// so preferable we would be keeping a goroutine alive per session/sub session
// Sessions can also be remote

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
