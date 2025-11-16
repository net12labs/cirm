package user

type Root struct {
	*User
}

func GetRoot() *Root {
	return &Root{
		User: NewUser(1, "root"),
	}
}

func (r *Root) UserCreateSys(id int64, name string) *User {
	// In a real implementation, you would have logic to generate unique IDs
	newUser := NewUser(id, name)
	return newUser
}

func (r *Root) UserCreate(name string) *User {
	// In a real implementation, you would have logic to generate unique IDs
	newUser := NewUser(1, name)
	return newUser
}
