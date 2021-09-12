package main

import (
    "fmt"

    "github.com/pkg/errors"
    "tcp/pkg/lib"
    "tcp/service"
)

func server() error {
    endpoint := lib.NewEndPoint()
    endpoint.AddHandleFunc("challenge", service.ChallengeHandler)
    endpoint.AddHandleFunc("verify", service.VerifyHandler)
    return endpoint.Listen()
}

func main() {
    err := server()
    if err != nil {
        fmt.Println("Error:", errors.WithStack(err))
    }
}
