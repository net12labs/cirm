package user

import "fmt"

type User struct {
	id       int64
	name     string
	globalId string
}

func NewUser(id int64, name string) *User {
	return &User{
		id:       id,
		name:     name,
		globalId: "89",
	}
}

func (u *User) ConsoleLog(item ...any) {
	fmt.Println(item...)
}

func (u *User) Id() int64 {
	return u.id
}

func (u *User) Name() string {
	return u.name
}
func (u *User) IsRoot() bool {
	return u.id == 1
}
func (u *User) IsNobody() bool {
	return u.id == 65534
}
func (u *User) IsSomebody() bool {
	return u.id == 65533
}

func (u *User) IsEverybody() bool {
	return u.id == 65532
}
func (u *User) IsAnybody() bool {
	return u.id == 65531
}

func (u *User) IsGuest() bool {
	return u.id == 65530
}

func (u *User) IsSystemUser() bool {
	return u.id < 1000
}

func (u *User) IsRegularUser() bool {
	return u.id >= 1000 && u.id != 65534
}

func (u *User) IsGroupMember(groupId int64) bool {
	// Placeholder implementation
	return false
}

func (u *User) GetGroups() []int64 {
	// Placeholder implementation
	return []int64{}
}
