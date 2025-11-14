package admin

import (
	"cirm/data/unit"
)

type Admin struct {
	Db *unit.SqliteDb
}

func NewAdmin() *Admin {
	return &Admin{}
}
