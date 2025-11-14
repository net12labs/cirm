package unit

type DataUnit struct {
	DB *SqliteDb
}

func NewDataUnit() *DataUnit {
	return &DataUnit{
		DB: NewDb(),
	}
}

func (du *DataUnit) Init() error {
	return du.DB.Init()
}
