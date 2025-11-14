package service

import (
	"github.com/net12labs/cirm/service-daemon/data/unit"
)

type Service struct {
	Db *unit.SqliteDb
}

func NewService() *Service {
	return &Service{}
}
