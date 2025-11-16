package cmd

// because the sender could be and agent posting to site api

type Cmd struct {
	Ttl       int
	ContextId int64
	SenderId  int64
	Cmd       string
	Src       []string
	Target    string // this should be a target context/target shell
	Args      []string
	Params    map[string]any
	Options   map[string]any
	Data      any
	StdOut    chan any
	StdErr    chan any
	StdIn     chan any
	Result    any
	ExitCode  int
	ErrorMsg  string
	// Command fields here
}

func NewCmd() *Cmd {
	return &Cmd{
		Params:   make(map[string]any),
		Options:  make(map[string]any),
		ExitCode: -1,
		Ttl:      5,
	}
}

func NewCmdWithChannels() *Cmd {
	return &Cmd{
		Params:   make(map[string]any),
		Options:  make(map[string]any),
		StdOut:   make(chan any),
		StdErr:   make(chan any),
		StdIn:    make(chan any),
		ExitCode: -1,
		Ttl:      5,
	}
}
