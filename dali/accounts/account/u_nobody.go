package account

type Nobody struct {
	*Account
}

func GetNobody() *Nobody {
	return &Nobody{
		Account: NewAccount(65534, "nobody"),
	}
}

type Anybody struct {
	*Account
}

func GetAnybody() *Anybody {
	return &Anybody{
		Account: NewAccount(65531, "anybody"),
	}
}

type Everybody struct {
	*Account
}

func GetEverybody() *Everybody {
	return &Everybody{
		Account: NewAccount(65532, "everybody"),
	}
}

type Somebody struct {
	*Account
}

func GetSomebody() *Somebody {
	return &Somebody{
		Account: NewAccount(65533, "somebody"),
	}
}

type SystemAccount struct {
	*Account
}

func GetSystemAccount(id int64, name string) *SystemAccount {
	return &SystemAccount{
		Account: NewAccount(id, name),
	}
}

type StdAccount struct {
	*Account
}

func GetStdAccount(id int64, name string) *StdAccount {
	return &StdAccount{
		Account: NewAccount(id, name),
	}
}

type GuestAccount struct {
	*Account
}

func GetGuestAccount() *GuestAccount {
	return &GuestAccount{
		Account: NewAccount(65530, "guest"),
	}
}
