package admin

import (
	"github.com/net12labs/cirm/dali/data"
)

type Admin struct {
	Db *data.SqliteDb
}

func NewAdmin() *Admin {
	return &Admin{}
}
