package cli

import (
	"log"

	"github.com/mitchellh/cli"
	"github.com/moorara/gocert/pki"
)

// App represents a cli app
type App struct {
	name    string
	version string
	init    cli.Command
	root    cli.Command
	interm  cli.Command
	server  cli.Command
	client  cli.Command
	sign    cli.Command
}

// NewApp creates a new cli app
func NewApp(name, version string) *App {
	return &App{
		name:    name,
		version: version,
		init:    NewInitCommand(),
		root:    NewReqCommand(pki.Metadata{CertType: pki.CertTypeRoot}),
		interm:  NewReqCommand(pki.Metadata{CertType: pki.CertTypeInterm}),
		server:  NewReqCommand(pki.Metadata{CertType: pki.CertTypeServer}),
		client:  NewReqCommand(pki.Metadata{CertType: pki.CertTypeClient}),
		sign:    NewSignCommand(),
	}
}

// Run executes the cli app
func (a *App) Run(args []string) int {
	app := cli.NewCLI(a.name, a.version)
	app.Args = args

	app.Commands = map[string]cli.CommandFactory{
		"init": func() (cli.Command, error) {
			return a.init, nil
		},
		"root": func() (cli.Command, error) {
			return a.root, nil
		},
		"intermediate": func() (cli.Command, error) {
			return a.interm, nil
		},
		"server": func() (cli.Command, error) {
			return a.server, nil
		},
		"client": func() (cli.Command, error) {
			return a.client, nil
		},
		"sign": func() (cli.Command, error) {
			return a.sign, nil
		},
	}

	status, err := app.Run()
	if err != nil {
		log.Println(err)
	}

	return status
}
