package cliserver

type CliServer struct {
	// CLI server fields here
	SocketPath string
}

func (cs *CliServer) Init() {
}
func (cs *CliServer) AddRoute(path string, handler func()) {
	// Add a new route to the CLI server
}
