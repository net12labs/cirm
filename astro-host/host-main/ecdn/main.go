package ecdn

import (
	"embed"

	ecdnwebserver "github.com/net12labs/cirm/mali/web-server-ecdn"
)

//go:embed cdn/*
var content embed.FS

type Ecdn struct {
	Server *ecdnwebserver.Server
}

func NewEcdn() *Ecdn {
	s := &Ecdn{Server: ecdnwebserver.NewServer()}
	s.Server.UrlBasePath = "/ecdn"
	s.Server.FsBasePath = "cdn"
	s.Server.Fs = &content
	return s
}
