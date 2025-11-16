package main

import (
	"fmt"

	"github.com/net12labs/cirm/service-daemon/cmd"
	service "github.com/net12labs/cirm/service-daemon/svc"
	"github.com/net12labs/cirm/service-daemon/unit"

	rtm "github.com/net12labs/cirm/dali/runtime"
)

// so context we can also package in the db - so it is all atomic

func main() {
	rtm.Etc.SetKV("unit_id", "default")
	rtm.Etc.SetKV("rtm_name", "china-ip-routes-maker")
	rtm.Etc.SetKV("pid_file_path", "../units/"+rtm.Etc.Get("unit_id").String()+"/proc/china-ip-routes-maker.pid")
	rtm.Etc.SetKV("data_dir", "../units/"+rtm.Etc.Get("unit_id").String()+"/data")
	rtm.Etc.SetKV("main_db_path", rtm.Etc.Get("data_dir").String()+"/main.db")

	rtm.Runtime.OnPanic.AddListener(func(err any) {
		fmt.Println("Runtime Panic:", err)
	})

	if rtm.Args.HasKey("--serve") {
		go func() {
			serve := service.NewServe()
			if err := serve.Start(); err != nil {
				fmt.Println("Failed to start service:", err)
				rtm.Runtime.Exit(1)
			}
			rtm.Runtime.Exit(0)
		}()
		rtm.Runtime.WaitForSIGTERM()
	}

	if rtm.Args.HasKey("--cmd") {
		cmd := cmd.NewCmd()
		if err := cmd.Execute(); err != nil {
			fmt.Println("Failed to execute command:", err)
			rtm.Runtime.Exit(1)
		}
		rtm.Runtime.Exit(0)
	}

	unit.PrintHelp()
	rtm.Runtime.Exit(0)

}
