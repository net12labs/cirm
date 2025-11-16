package shellcore

type ShellCore struct {
	AccountId int64
	Buffer    []string
	History   []string
}

func (s *ShellCore) AddToHistory(command string) {
	s.History = append(s.History, command)
}

func NewShellCore(accountId int64) *ShellCore {
	return &ShellCore{
		AccountId: accountId,
		Buffer:    []string{},
		History:   []string{},
	}
}

type Shell struct {
}

func NewShell() *Shell {
	return &Shell{}
}
