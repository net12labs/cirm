package main

import (
	"cirm/cmd"
	data "cirm/data"
	xmain "cirm/main"
	"fmt"
	"os"

	rtm "github.com/net12labs/cirm/dali/runtime"

	ox "github.com/net12labs/cirm/dali/ox"
)

// so context we can also package in the db - so it is all atomic

func main() {
	ox.Etc.SetKV("unit_id", "default")
	ox.Etc.SetKV("rtm_name", "china-ip-routes-maker")
	ox.Etc.SetKV("pid_file_name", "china-ip-routes-maker.pid")
	ox.Etc.SetKV("data_dir", "../../units")

	rtm.Runtime.OnPanic.AddListener(func(err any) {
		fmt.Println("Runtime Panic:", err)
	})

	fmt.Println("China IP Routes Maker", data.Module.DB.DbPath)

	if len(os.Args) > 1 && os.Args[1] == "--serve" {
		serve := xmain.NewServe()
		if err := serve.Start(); err != nil {
			fmt.Println("Failed to start service:", err)
			rtm.Runtime.Exit(1)
		}

	} else {
		cmd := cmd.Cmd{}
		cmd.OnExit = func() {
			// Cleanup tasks here
		}
		cmd.Execute(os.Args[1:])
		rtm.Runtime.OnExit.AddListener(func(code any) {
			fmt.Println("Exited with code at command", code)
		})
		rtm.Runtime.Exit(0)
	}

}
