package cli

import (
	"flag"

	"github.com/mitchellh/cli"
	"github.com/moorara/gocert/pki"
)

const (
	rootEnterConfig = "\nConfigurations for root certificate authority ..."
	rootEnterClaim  = "\nSpecifications for root certificate authority ..."

	rootSynopsis = ""
	rootHelp     = ``
)

// RootNewCommand represents the "root new" command
type RootNewCommand struct {
	ui  cli.Ui
	pki pki.Manager
}

// NewRootNewCommand creates a RootNewCommand
func NewRootNewCommand() *RootNewCommand {
	return &RootNewCommand{
		ui:  newColoredUI(),
		pki: pki.NewX509Manager(),
	}
}

// Synopsis returns the short help text for root command
func (c *RootNewCommand) Synopsis() string {
	return rootSynopsis
}

// Help returns the long help text for root command
func (c *RootNewCommand) Help() string {
	return rootHelp
}

// Run executes the root command
func (c *RootNewCommand) Run(args []string) int {
	flags := flag.NewFlagSet("root", flag.ContinueOnError)
	flags.Usage = func() {}
	err := flags.Parse(args)
	if err != nil {
		return ErrorInvalidFlag
	}

	state, spec, status := LoadWorkspace(c.ui)
	if status != 0 {
		return status
	}

	c.ui.Output(rootEnterConfig)
	AskForConfigCA(&state.Root, c.ui)

	c.ui.Output(rootEnterClaim)
	AskForClaim(&spec.Root, c.ui)

	err = c.pki.GenRootCA(state.Root, spec.Root)
	if err != nil {
		c.ui.Error("Failed to generate root ca. Error: " + err.Error())
		return ErrorRootCA
	}

	c.ui.Output("")
	return 0
}
