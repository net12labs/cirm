package meta

import "github.com/net12labs/cirm/ops/rtm"

func PrintHelp() {
	println("This is the help information for the Unit module.")
	println("Available commands:")
	println("  --serve   Start the service")
	println("  --cmd     Execute a command")
}

func Try() {
	if rtm.Args.HasKey("--help") || rtm.Args.HasKey("-h") {
		PrintHelp()
		rtm.Runtime.Exit(0)
	}
}
