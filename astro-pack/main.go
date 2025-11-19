package main

import (
	"github.com/net12labs/cirm/astro-pack/exec"
	proc_main "github.com/net12labs/cirm/astro-pack/main"
	"github.com/net12labs/cirm/astro-pack/meta"

	"github.com/net12labs/cirm/ops/rtm"

	config "github.com/net12labs/cirm/astro-pack/config"
)

func main() {

	config.Init()

	// so the domain needs to have separate logic (like a kernel)

	proc_main.Try()
	exec.Try()
	meta.Try()

	meta.PrintHelp()
	rtm.Runtime.Exit(0)

}
