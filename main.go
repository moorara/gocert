package main

import (
	"log"
	"os"

	"github.com/mitchellh/cli"
	"github.com/moorara/gocert/version"
)

var ()

func runApp(args []string) int {
	app := cli.NewCLI("gocert", version.GetFullSpec())
	app.Args = args
	app.Commands = map[string]cli.CommandFactory{}

	status, err := app.Run()
	if err != nil {
		log.Println(err)
	}

	return status
}

func main() {
	os.Exit(runApp(os.Args[1:]))
}
