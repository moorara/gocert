package main

import (
	"log"
	"os"

	"github.com/mitchellh/cli"
	"github.com/moorara/gotls/cmd"
	"github.com/moorara/gotls/version"
)

func runApp(args []string) int {
	app := cli.NewCLI("gotls", version.GetFullSpec())
	app.Args = args
	app.Commands = map[string]cli.CommandFactory{
		"cert": func() (cli.Command, error) {
			return cmd.NewCert(), nil
		},
	}

	status, err := app.Run()
	if err != nil {
		log.Println(err)
	}

	return status
}

func main() {
	os.Exit(runApp(os.Args[1:]))
}
