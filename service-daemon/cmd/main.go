package cmd

type Cmd struct {
	OnExit func()
	// Other root fields here
}

func (c *Cmd) Execute(args []string) {
	// Initialize other components here
	// Start the application
}
