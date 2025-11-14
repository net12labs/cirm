package data

import (
	"cirm/data/admin"
	"cirm/data/service"
	"cirm/data/unit"
)

var Module = unit.NewDataUnit()

var Service = service.NewService()
var Admin = admin.NewAdmin()

func init() {
	Service.Db = Module.DB
	Admin.Db = Module.DB
}
