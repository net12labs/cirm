package admin

import (
	"github.com/net12labs/cirm/service-daemon/data/unit"
)

type Admin struct {
	Db *unit.SqliteDb
}

func NewAdmin() *Admin {
	return &Admin{}
}
