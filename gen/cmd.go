package gen

import (
	"flag"

	"github.com/mitchellh/cli"
	"github.com/moorara/gocert/cert"
	"github.com/moorara/gocert/common"
	"github.com/moorara/gocert/config"
)

const (
	// ErrorInvalidFlags is returned when an invalid flag is provided
	ErrorInvalidFlags = 21
	// ErrorReadState is returned when cannot read state
	ErrorReadState = 22
	// ErrorReadSpec is returned when cannot read spec file
	ErrorReadSpec = 23
	// ErrorRootCA is returned when generating root ca failed
	ErrorRootCA = 24
	// ErrorIntermCA is returned when generating intermediate ca failed
	ErrorIntermCA = 25

	textSettingsEnterRoot   = "\nSettings for root certificate authority ..."
	textClaimEnterRoot      = "\nSpecs for root certificate authority ..."
	textSettingsEnterInterm = "\nSettings for intermediate certificate authority ..."
	textClaimEnterInterm    = "\nSpecs for intermediate certificate authority ..."

	textSynopsis = "Generates different types of certificates in a workspace"
	textHelp     = `
	You can use this command to generate the following types of certificates:
	  * Root Certificate Authority
	  * Intermediate Certificate Authority
	  * Server Certificate
	  * Client Certificate

	Flags:
	  -root
	  -intermediate
	  -server
	  -client
	`
)

// Command represents the gen command
type Command struct {
	ui   cli.Ui
	cert cert.Manager
}

// NewCommand creates a new command
func NewCommand() *Command {
	return &Command{
		ui:   common.NewColoredUI(),
		cert: cert.NewX509Manager(),
	}
}

// Synopsis returns the short help text for command
func (c *Command) Synopsis() string {
	return textSynopsis
}

// Help returns the long help text for command
func (c *Command) Help() string {
	return textHelp
}

// Run executes the command
func (c *Command) Run(args []string) int {
	var root, interm, server, client bool

	flags := flag.NewFlagSet("gen", flag.ContinueOnError)
	flags.Usage = func() {}
	flags.BoolVar(&root, "root", false, "generate the root certificate authority")
	flags.BoolVar(&interm, "intermediate", false, "generate an intermediate certificate authority")
	flags.BoolVar(&server, "server", false, "generate a server certificate")
	flags.BoolVar(&client, "client", false, "generate a client certificate")
	err := flags.Parse(args)
	if err != nil {
		return ErrorInvalidFlags
	}

	state, err := config.LoadState(config.FileNameState)
	if err != nil {
		c.ui.Error("Failed to read state " + config.FileNameState)
		return ErrorReadState
	}

	spec, err := config.LoadSpec(config.FileNameSpec)
	if err != nil {
		c.ui.Error("Failed to read spec " + config.FileNameSpec)
		return ErrorReadSpec
	}

	// Root CA
	if root {
		c.ui.Output(textSettingsEnterRoot)
		state.Root.FillIn(c.ui)

		c.ui.Output(textClaimEnterRoot)
		spec.Root.FillIn(c.ui)

		err := c.cert.GenRootCA(state.Root, spec.Root)
		if err != nil {
			c.ui.Error("Failed to generate root ca. Error: " + err.Error())
			return ErrorRootCA
		}
	}

	// Intermediate CA
	if interm {
		c.ui.Output(textSettingsEnterInterm)
		state.Interm.FillIn(c.ui)

		c.ui.Output(textClaimEnterInterm)
		spec.Interm.FillIn(c.ui)

		err := c.cert.GenIntermCA(state.Root, spec.Root)
		if err != nil {
			c.ui.Error("Failed to generate intermediate ca. Error: " + err.Error())
			return ErrorIntermCA
		}
	}

	c.ui.Output("")
	return 0
}
