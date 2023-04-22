package main

import (
	"os"

	"github.com/moorara/gocert/cli"
	"github.com/moorara/gocert/metadata"
)

func main() {
	app := cli.NewApp("gocert", metadata.String())
	status := app.Run(os.Args[1:])

	os.Exit(status)
}
