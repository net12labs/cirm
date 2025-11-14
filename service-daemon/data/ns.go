package data

import (
	"github.com/net12labs/cirm/service-daemon/data/admin"
	"github.com/net12labs/cirm/service-daemon/data/service"
	"github.com/net12labs/cirm/service-daemon/data/unit"
)

var Module = unit.NewDataUnit()

var Service = service.NewService()
var Admin = admin.NewAdmin()

func init() {
	Service.Db = Module.DB
	Admin.Db = Module.DB
}
