package data

import (
	mdata "github.com/net12labs/cirm/mali/data"
	mdo "github.com/net12labs/cirm/mali/rtm/do"
)

type ops struct {
	dbs map[string]*mdata.SqliteDb
	// internal fields
}

func Db(name string) *mdata.SqliteDb {
	return Ops.GetDb(name)
}

func (o *ops) CreateDb(name, path string) *mdata.SqliteDb {
	mdo.InitFsPath(path)

	o.dbs[name] = mdata.NewDb()
	o.dbs[name].DbPath = path
	return o.dbs[name]
}
func (o *ops) GetDb(name string) *mdata.SqliteDb {
	return o.dbs[name]
}

var Ops = &ops{
	dbs: make(map[string]*mdata.SqliteDb),
}
