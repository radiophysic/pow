package service

import (
	"log"

	"tcp/conf"
	"tcp/pkg/lib"
)

type Service struct {
	cfg conf.Config
	log *log.Logger
}

func NewService(config conf.Config, logger *log.Logger) Service {
	return Service{cfg: config, log: logger}
}

func (s Service) Server() error {
	endpoint := lib.NewEndPoint(s.log)
	endpoint.AddHandleFunc("challenge", s.ChallengeHandler)
	endpoint.AddHandleFunc("verify", s.VerifyHandler)
	return endpoint.Listen(s.cfg.ServerPort)
}
