package service

import (
	"tcp/conf"
	"tcp/pkg/lib"
)

type Service struct {
	cfg conf.Config
}

func NewService(config conf.Config) Service {
	return Service{cfg: config}
}

func (s Service) Server() error {
	endpoint := lib.NewEndPoint()
	endpoint.AddHandleFunc("challenge", s.ChallengeHandler)
	endpoint.AddHandleFunc("verify", s.VerifyHandler)
	return endpoint.Listen(s.cfg.ServerPort)
}
