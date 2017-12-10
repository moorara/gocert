package cli

import (
	"log"

	"github.com/mitchellh/cli"
	"github.com/moorara/gocert/pki"
)

// App represents a cli app
type App struct {
	name      string
	version   string
	init      cli.Command
	rootNew   cli.Command
	intermNew cli.Command
	serverNew cli.Command
	clientNew cli.Command
}

// NewApp creates a new cli app
func NewApp(name, version string) *App {
	return &App{
		name:      name,
		version:   version,
		init:      NewInitCommand(),
		rootNew:   NewRootNewCommand(),
		intermNew: NewReqCommand(pki.Metadata{CertType: pki.CertTypeInterm}),
		serverNew: NewReqCommand(pki.Metadata{CertType: pki.CertTypeServer}),
		clientNew: NewReqCommand(pki.Metadata{CertType: pki.CertTypeClient}),
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
		"root new": func() (cli.Command, error) {
			return a.rootNew, nil
		},
		"intermediate new": func() (cli.Command, error) {
			return a.intermNew, nil
		},
		"server new": func() (cli.Command, error) {
			return a.serverNew, nil
		},
		"client new": func() (cli.Command, error) {
			return a.clientNew, nil
		},
	}

	status, err := app.Run()
	if err != nil {
		log.Println(err)
	}

	return status
}
