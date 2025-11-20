package args_auth

type LoginCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (lc *LoginCredentials) IsValid() bool {
	return lc.Username != "" && lc.Password != ""
}
