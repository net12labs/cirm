package accounts

import (
	"github.com/net12labs/cirm/dali/accounts/account"
)

type Account = account.Account
type Group = account.Group

type accounts struct {
}

var Accounts = &accounts{}

func (u *accounts) AccountGetByName(name string) *account.Account {
	// In a real implementation, you would have logic to look up the Account by name
	// Here we just return a dummy Account for demonstration purposes
	return account.NewAccount(1, name)
}

// Also need to create the equivalent groups

var Root = account.GetRoot()
var Nobody = account.GetNobody()
var Somebody = account.GetSomebody()
var Anybody = account.GetAnybody()
var Everybody = account.GetEverybody()
var Guest = account.GetGuestAccount()
