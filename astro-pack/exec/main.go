package exec

import (
	"fmt"

	"github.com/net12labs/cirm/ops/rtm"
)

type Cmd struct {
	OnExit func()
	// Other root fields here
}

func (c *Cmd) Execute() error {
	// Initialize other components here
	// Start the application
	return nil
}

func NewCmd() *Cmd {
	cmd := &Cmd{}
	cmd.OnExit = func() {
		// Cleanup tasks here
	}
	return cmd
}

func Try() {
	if rtm.Args.HasKey("--exec") {
		cmd := NewCmd()
		if err := cmd.Execute(); err != nil {
			fmt.Println("Failed to execute command:", err)
			rtm.Runtime.Exit(1)
		}
		rtm.Runtime.Exit(0)
	}
}
