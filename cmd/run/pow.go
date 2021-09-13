package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/pkg/errors"
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
	flag.Parse()

	switch appType {
	case "server":
		err := Server()
		if err != nil {
			fmt.Println("Error:", errors.WithStack(err))
			os.Exit(1)
		}
	case "client":
		err := client(hostAddr)
		if err != nil {
			fmt.Println("Error:", errors.WithStack(err))
			os.Exit(1)
		}
	default:
		flag.Usage()
	}

	os.Exit(0)
}
