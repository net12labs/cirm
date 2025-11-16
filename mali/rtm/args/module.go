package args

import "os"

type Args struct {
	// Define argument fields here
}

func NewArgs() *Args {
	return &Args{}
}

func (a *Args) HasKey(key string) bool {
	if len(os.Args) <= 1 {
		return false
	}
	for i := 1; i < len(os.Args); i++ {
		if os.Args[i] == key {
			return true
		}
	}
	return false
}
