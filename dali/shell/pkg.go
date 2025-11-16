package shell

import (
	"github.com/net12labs/cirm/dali/bin"
	"github.com/net12labs/cirm/dali/users/user"
)

type ShellCore struct {
	User    *user.User
	Buffer  []string
	History []string
}

type Shell struct {
	Core *ShellCore
}

func NewShell(userId int64) *Shell {
	return &Shell{
		Core: &ShellCore{
			User:    user.NewUser(userId, ""),
			Buffer:  []string{},
			History: []string{},
		},
	}
}

func (s *Shell) Execute(path string, args ...string) int {
	err := bin.Exec(s.Core.User.Id(), path, args...)
	if err != nil {
		return 1
	}
	return 0
}

func (s *ShellCore) AddToHistory(command string) {
	s.History = append(s.History, command)
}

func (s *ShellCore) CreateStdUser(name string) *user.StdUser {
	return user.GetStdUser(5888, name)
}

// for every context - like root, project etc - there is a separate shell
// shell can also be shared - where we can have other users or AI agents

// the AI agent also works kind of like shell
// But we also need to be able to navigate the actual data - kind of like telnet into a user
