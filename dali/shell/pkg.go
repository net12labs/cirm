package shell

import (
	"github.com/net12labs/cirm/dali/accounts/account"
	"github.com/net12labs/cirm/dali/bin"
	shellcore "github.com/net12labs/cirm/mali/shell-core"
)

type Shell struct {
	*shellcore.Shell
	Core *shellcore.ShellCore
}

func NewShell(userId int64) *Shell {
	return &Shell{
		Shell: shellcore.NewShell(),
		Core:  shellcore.NewShellCore(userId),
	}
}

func (s *Shell) Execute(path string, args ...string) int {
	err := bin.Exec(s.Core.AccountId, path, args...)
	if err != nil {
		return 1
	}
	return 0
}

func (s *Shell) CreateStdAccount(name string) *account.StdAccount {
	return account.GetStdAccount(5888, name)
}
func (s *Shell) SetAccountPassword(accountId int64, password string) error {
	return nil
}

// for every context - like root, project etc - there is a separate shell
// shell can also be shared - where we can have other users or AI agents

// the AI agent also works kind of like shell
// But we also need to be able to navigate the actual data - kind of like telnet into a user
