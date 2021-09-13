package main

import (
	"tcp/pkg/lib"
	"tcp/service"
)

func Server() error {
	endpoint := lib.NewEndPoint()
	endpoint.AddHandleFunc("challenge", service.ChallengeHandler)
	endpoint.AddHandleFunc("verify", service.VerifyHandler)
	return endpoint.Listen()
}
