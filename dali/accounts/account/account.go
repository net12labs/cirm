package account

import "fmt"

type Account struct {
	id       int64
	name     string
	globalId string
}

func NewAccount(id int64, name string) *Account {
	return &Account{
		id:       id,
		name:     name,
		globalId: "89",
	}
}

func (u *Account) ConsoleLog(item ...any) {
	fmt.Println(item...)
}

func (u *Account) Id() int64 {
	return u.id
}

func (u *Account) Name() string {
	return u.name
}
func (u *Account) IsRoot() bool {
	return u.id == 1
}
func (u *Account) IsNobody() bool {
	return u.id == 65534
}
func (u *Account) IsSomebody() bool {
	return u.id == 65533
}

func (u *Account) IsEverybody() bool {
	return u.id == 65532
}
func (u *Account) IsAnybody() bool {
	return u.id == 65531
}

func (u *Account) IsGuest() bool {
	return u.id == 65530
}

func (u *Account) IsSystemAccount() bool {
	return u.id < 1000
}

func (u *Account) IsStdAccount() bool {
	return u.id >= 1000 && u.id != 65534
}

func (u *Account) IsGroupMember(groupId int64) bool {
	// Placeholder implementation
	return false
}

func (u *Account) GetGroups() []int64 {
	// Placeholder implementation
	return []int64{}
}
