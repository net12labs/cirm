package data

type ops struct {
	dbs map[string]*SqliteDb
	// internal fields
}

func Db(name string) *SqliteDb {
	return Ops.GetDb(name)
}

func (o *ops) CreateDb(name, path string) *SqliteDb {
	o.dbs[name] = NewDb()
	o.dbs[name].DbPath = path
	return o.dbs[name]
}
func (o *ops) GetDb(name string) *SqliteDb {
	return o.dbs[name]
}

var Ops = &ops{
	dbs: make(map[string]*SqliteDb),
}
