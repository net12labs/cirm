package exec

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
