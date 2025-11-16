package service

import (
	"github.com/net12labs/cirm/dali/data"
)

type Service struct {
	Db *data.SqliteDb
}

func NewService() *Service {
	return &Service{}
}
