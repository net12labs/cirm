package users

import (
	"github.com/net12labs/cirm/dali/users/user"
)

type User = user.User
type Group = user.Group

type users struct {
}

var Users = &users{}

func (u *users) UserGetByName(name string) *user.User {
	// In a real implementation, you would have logic to look up the user by name
	// Here we just return a dummy user for demonstration purposes
	return user.NewUser(1, name)
}

var Root = user.GetRoot()
var Nobody = user.GetNobody()
var Somebody = user.GetSomebody()
var Anybody = user.GetAnybody()
var Everybody = user.GetEverybody()
var Guest = user.GetGuestUser()
