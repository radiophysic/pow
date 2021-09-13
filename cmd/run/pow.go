package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/pkg/errors"
	"tcp/client"
	"tcp/conf"
	"tcp/service"
)

var (
	appType  string
	hostAddr string
)

func init() {
	flag.StringVar(&appType, "mode", "", "-mode=server\n-app=client")
	flag.StringVar(&hostAddr, "host", "pow_server", "-host=localhost")
}

func main() {
	var (
		cfg conf.Config
		err error
		svc service.Service
		cli client.Client
	)

	flag.Parse()

	if cfg, err = conf.ReadConfig(); err != nil {
		fmt.Println("Error:", errors.WithStack(err))
		os.Exit(1)
	}

	switch appType {
	case "server":
		svc = service.NewService(cfg)
		err = svc.Server()
		if err != nil {
			fmt.Println("Error:", errors.WithStack(err))
			os.Exit(1)
		}
	case "client":
		cli = client.NewClient(cfg)
		err = cli.GetQuote()
		if err != nil {
			fmt.Println("Error:", errors.WithStack(err))
			os.Exit(1)
		}
	default:
		flag.Usage()
	}

	os.Exit(0)
}
