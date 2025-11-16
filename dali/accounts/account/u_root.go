package account

type Root struct {
	*Account
}

func GetRoot() *Root {
	return &Root{
		Account: NewAccount(1, "root"),
	}
}

func (r *Root) AccountCreateSys(id int64, name string) *Account {
	// In a real implementation, you would have logic to generate unique IDs
	newAccount := NewAccount(id, name)
	return newAccount
}

func (r *Root) AccountCreate(name string) *Account {
	// In a real implementation, you would have logic to generate unique IDs
	newAccount := NewAccount(1, name)
	return newAccount
}
