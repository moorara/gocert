package main

import (
	"os"

	"github.com/moorara/gocert/cli"
	"github.com/moorara/gocert/version"
)

func main() {
	app := cli.NewApp("gocert", version.String())
	status := app.Run(os.Args[1:])

	os.Exit(status)
}
