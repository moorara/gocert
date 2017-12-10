package cli

import (
	"flag"

	"github.com/mitchellh/cli"
	"github.com/moorara/gocert/pki"
)

const (
	rootName = "root"

	rootEnterConfig = "\nCONFIGURATIONS FOR ROOT CERTIFICATE AUTHORITY ..."
	rootEnterClaim  = "\nSPECIFICATIONS FOR ROOT CERTIFICATE AUTHORITY ..."

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
	flags := flag.NewFlagSet("root new", flag.ContinueOnError)
	flags.Usage = func() {}
	err := flags.Parse(args)
	if err != nil {
		return ErrorInvalidFlag
	}

	state, spec, status := loadWorkspace(c.ui)
	if status != 0 {
		return status
	}

	md := pki.Metadata{CertType: pki.CertTypeRoot}

	c.ui.Output(rootEnterConfig)
	askForConfig(&state.Root, c.ui)

	c.ui.Output(rootEnterClaim)
	askForClaim(&spec.Root, c.ui)

	c.ui.Output("")

	err = c.pki.GenCert(rootName, state.Root, spec.Root, md)
	if err != nil {
		c.ui.Error("Failed to generate root ca. Error: " + err.Error())
		return ErrorRootCA
	}

	return 0
}
