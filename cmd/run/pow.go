package main

import (
	"flag"
	"log"
	"os"

	"github.com/pkg/errors"
	"tcp/client"
	"tcp/conf"
	"tcp/service"
)

func main() {
	var (
		cfg conf.Config
		err error
		svc service.Service
		cli client.Client
		appType  string
		hostAddr string
	)

	flag.StringVar(&appType, "mode", "", "-mode=server\n-app=client")
	flag.StringVar(&hostAddr, "host", "pow_server", "-host=localhost")
	flag.Parse()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	if cfg, err = conf.ReadConfig(); err != nil {
		logger.Println("Error:", errors.WithStack(err))
		os.Exit(1)
	}

	switch appType {
	case "server":
		svc = service.NewService(cfg, logger)
		err = svc.Server()
		if err != nil {
			logger.Println("Error:", errors.WithStack(err))
			os.Exit(1)
		}
	case "client":
		cli = client.NewClient(cfg, logger)
		err = cli.GetQuote()
		if err != nil {
			logger.Println("Error:", errors.WithStack(err))
			os.Exit(1)
		}
	default:
		flag.Usage()
	}

	os.Exit(0)
}
