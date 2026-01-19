package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
)

const help = `Usage: go run ./cmd/build/main.go [OPTIONS] COMMAND

CLI tools for weaviate-go-client maintainers.

Commands:
	contracts	Update ./api/proto/ and ./api/rest specs to their latest versions in github.com/weaviate/weaviate.

Run go run ./cmd/build COMMAND --help for more information on the command.

Global Options:
	-h, --help	Print this message.
`

func main() {
	ctx := context.Background()

	h := flag.Bool("help", false, "Print help message.")
	flag.Parse()

	if *h {
		fmt.Print(help)
		os.Exit(0)
	}

	if len(os.Args) < 2 {
		fmt.Println("ERROR: not enough arguments")
		fmt.Print(help)
		os.Exit(1)
	}

	var err error
	switch os.Args[1] {
	case "contracts":
		err = Contracts(ctx, os.Args[2:])
	default:
		err = errors.New(help)
	}

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
