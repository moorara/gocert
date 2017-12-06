package main

import (
	"log"
	"os"

	"github.com/mitchellh/cli"
	"github.com/moorara/gocert/gen"
	"github.com/moorara/gocert/new"
	"github.com/moorara/gocert/version"
)

var (
	cmdNew cli.Command = new.NewCommand()
	cmdGen cli.Command = gen.NewCommand()
)

func runApp(args []string) int {
	app := cli.NewCLI("gocert", version.GetFullSpec())
	app.Args = args
	app.Commands = map[string]cli.CommandFactory{
		"new": func() (cli.Command, error) {
			return cmdNew, nil
		},
		"gen": func() (cli.Command, error) {
			return cmdGen, nil
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
