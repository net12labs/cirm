package user

type Nobody struct {
	*User
}

func GetNobody() *Nobody {
	return &Nobody{
		User: NewUser(65534, "nobody"),
	}
}

type Anybody struct {
	*User
}

func GetAnybody() *Anybody {
	return &Anybody{
		User: NewUser(65531, "anybody"),
	}
}

type Everybody struct {
	*User
}

func GetEverybody() *Everybody {
	return &Everybody{
		User: NewUser(65532, "everybody"),
	}
}

type Somebody struct {
	*User
}

func GetSomebody() *Somebody {
	return &Somebody{
		User: NewUser(65533, "somebody"),
	}
}

type SystemUser struct {
	*User
}

func GetSystemUser(id int64, name string) *SystemUser {
	return &SystemUser{
		User: NewUser(id, name),
	}
}

type StdUser struct {
	*User
}

func GetStdUser(id int64, name string) *StdUser {
	return &StdUser{
		User: NewUser(id, name),
	}
}

type GuestUser struct {
	*User
}

func GetGuestUser() *GuestUser {
	return &GuestUser{
		User: NewUser(65530, "guest"),
	}
}
