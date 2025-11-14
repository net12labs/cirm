package service

import (
	"cirm/data/unit"
)

type Service struct {
	Db *unit.SqliteDb
}

func NewService() *Service {
	return &Service{}
}
